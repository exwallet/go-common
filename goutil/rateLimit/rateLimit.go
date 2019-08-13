/*
 * @Author: kidd
 * @Date: 7/30/19 10:52 AM
 */

package rateLimit

import (
	"github.com/exwallet/go-common/goutil/gotime"
	"sync"
)

/**
--------------------------------------------------------------------------------------------------> tryTimes
       pass             |         fail and wait         |         fail and wait and punish
                        |                               |
                duMaxTryTimes                 duMaxFailTimes


*/
type rateLimit struct {
	duCount             int64           // 周期数量
	duUnitMillSec       int64           // 周期时长,微秒
	duMaxTryTimes       int64           // 周期内允许尝试次数
	duMaxFailTimes      int64           // 周期内允许最大失败次数
	tryTimesMap         map[int64]int64 // 周期热点映射表; map[某个周期]次数, tryTimesMap[0]:当前周期内, tryTimesMap[1]:一个周期前
	activeDuUnitMillSec int64           //
	activeDuMaxTryTimes int64           //
	punishFactor        float64         // (0<n<1)惩罚倍数因子,失败次数达多少时, 周期变长, 允许次数变小 .
	lastPeriodTime      int64           // 最新周期时间点
	hasPunish           bool            //
	//
	lock sync.Mutex //
}

type durationCore struct {
	tryTimes int64
	punish   bool
}

func NewRateLimit(durationUnitMillSec int64, durationCount int64, durationMaxTryTime int64, durationMaxFailTimes int64, punishFactor... float64) *rateLimit {
	if durationUnitMillSec <= 0 || durationCount <=0 || durationMaxTryTime <= 0 || durationMaxFailTimes < durationMaxTryTime {
		panic("RateLimit非法参数")
	}

	a := &rateLimit{
		duCount:             durationCount,
		duUnitMillSec:       durationUnitMillSec,
		duMaxTryTimes:       durationMaxTryTime,
		duMaxFailTimes:      durationMaxFailTimes,
		tryTimesMap:         make(map[int64]int64, durationCount),
		activeDuUnitMillSec: durationUnitMillSec,
		activeDuMaxTryTimes: durationMaxTryTime,
		punishFactor:        0,
		lastPeriodTime:      gotime.UnixNowMillSec(),
		hasPunish:           false,
		lock:                sync.Mutex{},
	}
	if len(punishFactor) > 0{
		if punishFactor[0] < 0 || punishFactor[0] >= 1 {
			panic("RateLimit 惩罚因子非法定义")
		}
		a.punishFactor = punishFactor[0]
	}
	for i := int64(0); i < durationCount; i++ {
		a.tryTimesMap[i] = 0
	}
	return a
}

func (r *rateLimit) rotate(now int64) {
	if r.tryTimesMap[0] < r.duMaxFailTimes {
		// 重置
		r.activeDuUnitMillSec = r.duUnitMillSec
		r.activeDuMaxTryTimes = r.duMaxTryTimes
	}
	r.hasPunish = false
	//
	for i := r.duCount - 1; i > 0; i-- {
		r.tryTimesMap[i] = r.tryTimesMap[i-1]
	}
	r.tryTimesMap[0] = 1
	r.lastPeriodTime = now
}

// 惩罚机制: 周期时间延时, 周期允许尝试次数缩小
func (r *rateLimit) doPunish() {
	r.activeDuUnitMillSec += int64(float64(r.activeDuUnitMillSec) * r.punishFactor)
	r.activeDuMaxTryTimes -= int64(float64(r.duMaxTryTimes) * r.punishFactor)
	if r.activeDuMaxTryTimes <= 0 {
		r.activeDuMaxTryTimes = 1
	}
	r.hasPunish = true
}

// 不通过返回等待恢复时间, 秒
func (r *rateLimit) Pass() (pass bool, coolSec int64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	now := gotime.UnixNowMillSec()
	tsDiff := now - r.lastPeriodTime
	if tsDiff <= r.activeDuUnitMillSec {
		r.tryTimesMap[0] += 1
		t := r.tryTimesMap[0]
		if t <= r.duMaxTryTimes {
			// do pass
			return true, 0
		}
		// do fail
		if t <= r.duMaxFailTimes {
			coolSec = (r.activeDuUnitMillSec - tsDiff)/1000
			return false, coolSec
		}
		// do fail and punish
		if !r.hasPunish {
			r.doPunish()
		}
		coolSec = (r.activeDuUnitMillSec - tsDiff)/1000
		return false, coolSec
	}
	// do pass
	r.rotate(now)
	return true, 0
}

// 取得周期热点列表
func (r *rateLimit) GetTryMap() []int64 {
	var out = make([]int64, r.duCount)
	for i, v := range r.tryTimesMap {
		out[i] = v
	}
	return out
}
