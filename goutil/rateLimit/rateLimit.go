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
                durationMaxTryTimes                 durationMaxFailTimes


*/
type rateLimit struct {
	durationCount             int64           // 周期数量
	durationUnitSec           int64           // 周期时长,秒
	durationMaxTryTimes       int64           // 周期内允许尝试次数
	durationMaxFailTimes      int64           // 周期内允许最大失败次数
	tryTimesMap               map[int64]int64 // 周期热点映射表; map[某个周期]次数, tryTimesMap[0]:当前周期内, tryTimesMap[1]:一个周期前
	activeDurationUnitSec     int64           //
	activeDurationMaxTryTimes int64           //
	punishFactor              float64         // (0<n<1)惩罚倍数因子,失败次数达多少时, 周期变长, 允许次数变小 .
	lastPeriodTime            int64           // 最新周期时间点
	hasPunish                 bool            //
	//
	lock sync.Mutex //
}

type durationCore struct {
	tryTimes int64
	punish   bool
}

func NewRateLimit(durationUnitSec int64, durationCount int64, durationMaxTryTime int64, durationMaxFailTimes int64, punishFactor float64) *rateLimit {
	a := &rateLimit{
		durationCount:             durationCount,
		durationUnitSec:           durationUnitSec,
		durationMaxTryTimes:       durationMaxTryTime,
		durationMaxFailTimes:      durationMaxFailTimes,
		tryTimesMap:               make(map[int64]int64, durationCount),
		activeDurationUnitSec:     durationUnitSec,
		activeDurationMaxTryTimes: durationMaxTryTime,
		punishFactor:              punishFactor,
		lastPeriodTime:            gotime.UnixNowSec(),
		hasPunish:                 false,
		lock:                      sync.Mutex{},
	}
	for i := int64(0); i < durationCount; i++ {
		a.tryTimesMap[i] = 0
	}
	return a
}

func (r *rateLimit) rotate(now int64) {
	r.hasPunish = false
	for i := r.durationCount - 1; i > 0; i-- {
		r.tryTimesMap[i] = r.tryTimesMap[i-1]
	}
	r.tryTimesMap[0] = 1
	r.lastPeriodTime = now
	// 恢复默认值
	if r.durationCount > 1 && r.tryTimesMap[1] < r.durationMaxFailTimes {
		r.activeDurationUnitSec = r.durationUnitSec
		r.activeDurationMaxTryTimes = r.durationMaxTryTimes
	}
}

// 惩罚机制: 周期时间延时, 周期允许尝试次数缩小
func (r *rateLimit) doPunish() {
	r.activeDurationUnitSec += int64(float64(r.activeDurationUnitSec) * r.punishFactor)
	r.activeDurationMaxTryTimes -= int64(float64(r.durationMaxTryTimes) * r.punishFactor)
	if r.activeDurationMaxTryTimes <= 0 {
		r.activeDurationMaxTryTimes = 1
	}
	r.hasPunish = true
}

// 不通过返回等待恢复时间, 秒
func (r *rateLimit) Pass() (pass bool, coolSec int64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	now := gotime.UnixNowSec()
	tsDiff := now - r.lastPeriodTime
	if tsDiff <= r.activeDurationUnitSec {
		r.tryTimesMap[0] += 1
		t := r.tryTimesMap[0]
		if t <= r.durationMaxTryTimes {
			// do pass
			return true, 0
		}
		// do fail
		if t <= r.durationMaxFailTimes {
			coolSec = r.activeDurationUnitSec - tsDiff
			return false, coolSec
		}
		// do fail and punish
		if !r.hasPunish {
			r.doPunish()
		}
		coolSec = r.activeDurationUnitSec - tsDiff
		return false, coolSec
	}
	// do pass
	r.rotate(now)
	return true, 0

}

// 取得周期热点列表
func (r *rateLimit) GetTryMap() []int64 {
	var out = make([]int64, r.durationCount)
	for i, v := range r.tryTimesMap {
		out[i] = v
	}
	return out
}
