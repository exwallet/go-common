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
                DuMaxTryTimes                 DuMaxFailTimes


*/
type rateLimit struct {
	DuCount             int64           // 周期数量
	DuUnitSec           int64           // 周期时长,微秒
	DuMaxTryTimes       int64           // 周期内允许尝试次数
	DuMaxFailTimes      int64           // 周期内允许最大失败次数
	PunishFactor        float64         // (0<n<1)惩罚倍数因子,失败次数达多少时, 周期变长, 允许次数变小 .
	tryTimesMap         map[int64]int64 // 周期热点映射表; map[某个周期]次数, tryTimesMap[0]:当前周期内, tryTimesMap[1]:一个周期前
	activeDuUnitSec     int64           //
	activeDuMaxTryTimes int64           //
	lastPeriodTime      int64           // 最新周期时间点
	hasPunish           bool            //
	//
	lock sync.Mutex //
}

type durationCore struct {
	tryTimes int64
	punish   bool
}

// durationMaxFailTimes <= 0 不惩罚
func NewRateLimit(durationUnitSec int64, durationCount int64, durationMaxTryTime int64, durationMaxFailTimes int64, punishFactor... float64) *rateLimit {
	if durationUnitSec <= 0 || durationCount <=0 || durationMaxTryTime <= 0 {
		panic("RateLimit非法参数")
	}

	a := &rateLimit{
		DuCount:             durationCount,
		DuUnitSec:           durationUnitSec,
		DuMaxTryTimes:       durationMaxTryTime,
		DuMaxFailTimes:      durationMaxFailTimes,
		PunishFactor:        0,
		tryTimesMap:         make(map[int64]int64, durationCount),
		activeDuUnitSec:     durationUnitSec,
		activeDuMaxTryTimes: durationMaxTryTime,
		lastPeriodTime:      gotime.UnixNowSec(),
		hasPunish:           false,
		lock:                sync.Mutex{},
	}
	if len(punishFactor) > 0{
		if punishFactor[0] < 0 || punishFactor[0] >= 1 {
			panic("RateLimit 惩罚因子非法定义")
		}
		a.PunishFactor = punishFactor[0]
	}
	for i := int64(0); i < durationCount; i++ {
		a.tryTimesMap[i] = 0
	}
	return a
}

func (r *rateLimit) rotate(now int64) {
	if r.tryTimesMap[0] < r.DuMaxFailTimes {
		// 重置
		r.activeDuUnitSec = r.DuUnitSec
		r.activeDuMaxTryTimes = r.DuMaxTryTimes
	}
	r.hasPunish = false
	//
	for i := r.DuCount - 1; i > 0; i-- {
		r.tryTimesMap[i] = r.tryTimesMap[i-1]
	}
	r.tryTimesMap[0] = 1
	r.lastPeriodTime = now
}

// 惩罚机制: 周期时间延时, 周期允许尝试次数缩小
func (r *rateLimit) doPunish() {
	r.activeDuUnitSec += int64(float64(r.activeDuUnitSec) * r.PunishFactor)
	r.activeDuMaxTryTimes -= int64(float64(r.DuMaxTryTimes) * r.PunishFactor)
	if r.activeDuMaxTryTimes <= 0 {
		r.activeDuMaxTryTimes = 1
	}
	r.hasPunish = true
}

// 不通过返回等待恢复时间, 秒
func (r *rateLimit) Pass() (pass bool, coolSec int64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	now := gotime.UnixNowSec()
	tsDiff := now - r.lastPeriodTime
	if tsDiff <= r.activeDuUnitSec {
		r.tryTimesMap[0] += 1
		t := r.tryTimesMap[0]
		if t <= r.DuMaxTryTimes {
			// do pass
			return true, 0
		}
		// do fail
		if t <= r.DuMaxFailTimes {
			coolSec = r.activeDuUnitSec - tsDiff
			return false, coolSec
		}
		// do fail and punish
		if r.DuMaxFailTimes <= 0 && !r.hasPunish {
			r.doPunish()
		}
		coolSec = r.activeDuUnitSec - tsDiff
		return false, coolSec
	}
	// do pass
	r.rotate(now)
	return true, 0
}

// 取得周期热点列表
func (r *rateLimit) GetTryMap() []int64 {
	var out = make([]int64, r.DuCount)
	for i, v := range r.tryTimesMap {
		out[i] = v
	}
	return out
}
