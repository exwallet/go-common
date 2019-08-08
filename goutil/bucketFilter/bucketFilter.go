/*
 * @Author: kidd
 * @Date: 1/24/19 12:15 PM
 *
 */

package bucketFilter

import (
	"errors"
	"github.com/exwallet/go-common/cache/redis"
	"github.com/exwallet/go-common/gologger"
	"github.com/exwallet/go-common/goutil/gotime"
)

/**
基于时间窗口的频率限制器, 最小单位秒
*/

var log = gologger.GetLogger()

// 过滤桶,限制发送动作
type BucketFilter struct {
	Running                   bool   // 正在运行
	Name                      string // 名称
	CacheKeyPrefix            string // 缓存前缀
	PeriodSeconds             int64  // 单位时间
	PeriodMaxTimes            int    // 单位时间内允许总次数
	AllowRetryIntervalSeconds int64  // 间隔多少秒后允许重试
	Stat                      Stat   // 统计
	ErrorRetryTooFast         string // 重试频繁错误提示
	ErrorTryTooMuch           string // 尝试次数过多错误提示
}

type Stat struct {
	StartTime   int64 // 启动时间
	BlockCount  int64 // 拦截总次数
	PermitCount int64 // 允许总次数
}

//var RetryTooFastError = errors.New("请求频率过高, 稍等一下") // 重次频繁
//var TryTooMuchError = errors.New("尝试次数过多,请休息一下")   // 尝试次数超过总允许次数

var ActiveBucketFilers []*BucketFilter

func NewBucketFilter(name string, CacheKeyPrefix string, PeriodSeconds int64, PeriodMaxTimes int, AllowRetryIntervalSeconds int64,
	errorRetryTooFastAndErrorTryTooMuch ...string) *BucketFilter {
	b := &BucketFilter{
		Running:                   true,
		Name:                      name,
		CacheKeyPrefix:            CacheKeyPrefix,
		PeriodSeconds:             PeriodSeconds,
		PeriodMaxTimes:            PeriodMaxTimes,
		AllowRetryIntervalSeconds: AllowRetryIntervalSeconds,
		Stat: Stat{
			StartTime:   gotime.UnixNow(),
			BlockCount:  0,
			PermitCount: 0,
		},
		ErrorRetryTooFast: "",
		ErrorTryTooMuch:   "",
	}
	if len(errorRetryTooFastAndErrorTryTooMuch) > 1 {
		b.ErrorRetryTooFast = errorRetryTooFastAndErrorTryTooMuch[0]
		b.ErrorTryTooMuch = errorRetryTooFastAndErrorTryTooMuch[1]
	} else {
		b.ErrorRetryTooFast = "请求频率过高, 稍等一下"
		b.ErrorTryTooMuch = "尝试次数过多,请休息一下"
	}
	ActiveBucketFilers = append(ActiveBucketFilers, b)
	return b
}

type cacheData struct {
	D []int64
}

func (b *BucketFilter) getCacheKey(key string) string {
	return b.CacheKeyPrefix + key
}

// 返回有效单位时间的时间窗口
func (b *BucketFilter) getTimeWindow(now int64, key string) ([]int64, error) {
	obj, e := redis.GetObj(b.getCacheKey(key), (*cacheData)(nil))
	if e != nil {
		return nil, e
	}
	if obj == nil {
		return []int64{}, nil
	}
	tw, ok := obj.(*cacheData)
	if !ok {
		return nil, errors.New("调用失败 (0091)")
	}
	if len(tw.D) == 0 {
		return []int64{}, nil
	}
	// 如果窗口最早的时间点在统计时间范围内, 直接返回
	if now-tw.D[0] < b.PeriodSeconds {
		gologger.Debug("------> 当前时间 %d\n", gotime.UnixNow())
		gologger.Debug("------> 返回原时间窗口 %+v\n", tw.D)
		return tw.D, nil
	}
	// 重新组装窗口
	var newTw []int64
	for _, t := range tw.D {
		if now-t < b.PeriodSeconds {
			newTw = append(newTw, t)
		}
	}
	gologger.Debug("------> 当前时间 %d\n", gotime.UnixNow())
	gologger.Debug("------> 重装时间窗口 %+v\n", newTw)
	return newTw, nil
}

func (b *BucketFilter) Stop() {
	b.Running = false
}

func (b *BucketFilter) Start() {
	b.Running = true
}

// 检查是否允许通过
func (b *BucketFilter) CheckAvailable(key string) error {
	if !b.Running {
		return nil
	}
	now := gotime.UnixNow()
	window, e := b.getTimeWindow(now, key)
	if e != nil {
		b.Stat.BlockCount += 1
		log.Error("缓存出错: %s", e.Error())
		return e
	}
	width := len(window)
	if width == 0 {
		return nil
	}
	// 检查最后次数与当前时间的间隔
	if now-window[width-1] < b.AllowRetryIntervalSeconds {
		b.Stat.BlockCount += 1
		log.Warn("拦截命中:%s %s %s", b.Name, b.ErrorRetryTooFast, key)
		return errors.New(b.ErrorRetryTooFast)
	}
	// 检查有效时间窗口内的尝试次数
	if len(window) >= b.PeriodMaxTimes {
		b.Stat.BlockCount += 1
		log.Warn("拦截命中:%s %s %s", b.Name, b.ErrorTryTooMuch, key)
		return errors.New(b.ErrorTryTooMuch)
	}
	return nil
}

// 成功触发更新缓存的时间窗口
func (b *BucketFilter) Success(key string) {
	if !b.Running {
		return
	}
	now := gotime.UnixNow()
	w, e := b.getTimeWindow(now, key)
	var d *cacheData
	if e == nil {
		window := append(w, now)
		d = &cacheData{D: window}
	} else {
		d = &cacheData{D: []int64{now}}
	}
	redis.SetObjAndExpire(b.getCacheKey(key), d, int(b.PeriodSeconds))
	b.Stat.PermitCount += 1
}
