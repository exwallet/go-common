/*
 * @Author: kidd
 * @Date: 7/30/19 10:52 AM
 */

package rateLimit

import (
	"github.com/exwallet/go-common/goutil/gotime"
	"sync"
	"time"
)

type rateLimit2 struct {
	periodCount    int64           // 周期数量
	periodUnit     int64           // 周期时长,秒
	periodMaxTimes int64           // 周期内允许次数
	hitsMap        map[int64]int64 // 周期热点映射表; map[某个周期]次数, hitsMap[0]:当前周期内, hitsMap[1]:一个周期前
	failsMap       map[int64]int64 // 失败次数
	lastTime       int64           // 最新周期时间点
	// todo:  惩罚因子, 连续失败次数达多少时, 周期变长, 允许次数变小 .
	lock           sync.Mutex      //
}

func NewRateLimit2(periodUnitSec int64, periodCount int64, periodMaxTimes int64) *rateLimit2 {
	a := &rateLimit2{
		periodCount:    periodUnitSec,
		periodUnit:     periodCount,
		periodMaxTimes: periodMaxTimes,
		lastTime:       gotime.UnixNowSec(),
		hitsMap:        make(map[int64]int64, periodCount),
		failsMap:       make(map[int64]int64, periodCount),
	}
	for i := int64(0); i < periodCount; i++ {
		a.hitsMap[i] = 0
		a.failsMap[i] = 0
	}
	return a
}

// 不通过返回等待恢复时间, 秒
func (r *rateLimit2) Pass() (pass bool, coolSec int64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	c := r.hitsMap[0]
	//
	now := gotime.UnixNowSec()
	if now-r.lastTime <= r.periodUnit {
		// 在最后有效周期内, 判断次数
		if c+1 > r.periodMaxTimes {
			coolSec = r.periodUnit - (now - r.lastTime)
			r.failsMap[0] = r.failsMap[0] + 1
			return false, coolSec
		}
		r.hitsMap[0] = c + 1
		return true, 0
	}
	// 各周期顺位下移
	for i := r.periodCount - 1; i > 0; i-- {
		r.hitsMap[i] = r.hitsMap[i-1]
		r.failsMap[i] = r.failsMap[i-1]
	}
	r.hitsMap[0] = 1
	r.failsMap[0] = 0
	r.lastTime = now
	return true, 0
}

// 取得周期热点列表
func (r *rateLimit2) GetHitsMap() []int64 {
	var out = make([]int64, r.periodCount)
	for i, v := range r.hitsMap {
		out[i] = v
	}
	return out
}
// 取得失败次数列表
func (r *rateLimit2) GetFailsMap() []int64 {
	var out = make([]int64, r.periodCount)
	for i, v := range r.failsMap {
		out[i] = v
	}
	return out
}

// =============================================================

/*
基于固定时间差内的频率控制, 时间差越大,误差也大, 适合短时间内高频调用限速
*/

// 限制速率
func NewRateLimit4Wait(ratePerSec int) *rateLimit4Wait {
	return &rateLimit4Wait{
		interval: time.Microsecond * time.Duration(1000*1000/ratePerSec),
		begin:    time.Now(),
	}
}

type rateLimit4Wait struct {
	interval time.Duration
	begin    time.Time
	lock     sync.Mutex
}

func (r *rateLimit4Wait) Pass() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if time.Now().Sub(r.begin) < r.interval {
		time.Sleep(r.interval)
	}
	r.begin = time.Now()
	return
}

// =================================================

func NewRateLimit(rate int, periodSec int) *rateLimit {
	return &rateLimit{
		rate:     rate,
		interval: time.Second * time.Duration(periodSec),
		begin:    time.Now(),
		count:    0,
		lock:     sync.Mutex{},
	}
}

type rateLimit struct {
	rate     int
	interval time.Duration
	begin    time.Time
	count    int
	lock     sync.Mutex
}

func (r *rateLimit) Pass() bool {
	result := true
	r.lock.Lock()
	defer r.lock.Unlock()
	//达到每秒速率限制数量
	//大于则速率在允许范围内，开始重新记数，返回true
	//小于，则返回false，记数不变
	if r.count == r.rate {
		if time.Now().Sub(r.begin) >= r.interval {
			//速度允许范围内，开始重新记数
			r.begin = time.Now()
			r.count = 0
		} else {
			result = false
		}
	} else {
		//没有达到速率限制数量，记数加1
		r.count++
	}
	return result
}
