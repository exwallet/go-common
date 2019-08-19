/*
 * @Author: kidd
 * @Date: 7/30/19 10:52 AM
 */

package limiter

import (
	"github.com/exwallet/go-common/goutil/gotime"
)

/**
限速器
--------------------------------------------------------------------------------------------------> tryTimes
       pass             |         fail and wait         |         fail and wait and punish
                        |                               |
                DuMaxTryTimes                 DuMaxFailTimes


*/
type RateLimiter struct {
	DuCount             int64           `json:"duCount"`             // 周期数量
	DuUnitSec           int64           `json:"duUnitSec"`           // 周期时长,微秒
	DuMaxTryTimes       int64           `json:"duMaxTryTimes"`       // 周期内允许尝试次数
	DuMaxFailTimes      int64           `json:"duMaxFailTimes"`      // 周期内允许最大失败次数
	PunishFactor        float64         `json:"punishFactor"`        // (0<n<1)惩罚倍数因子,失败次数达多少时, 周期变长, 允许次数变小 .
	TryTimesMap         map[int64]int64 `json:"tryTimesMap"`         // 周期热点映射表; map[某个周期]次数, TryTimesMap[0]:当前周期内, TryTimesMap[1]:一个周期前
	ActiveDuUnitSec     int64           `json:"activeDuUnitSec"`     //
	ActiveDuMaxTryTimes int64           `json:"activeDuMaxTryTimes"` //
	LastPeriodTime      int64           `json:"lastPeriodTime"`      // 最新周期时间点
	HasPunish           bool            `json:"hasPunish"`           //
}

type durationCore struct {
	tryTimes int64
	punish   bool
}

// durationMaxFailTimes <= 0 不惩罚
func NewRateLimiter(durationUnitSec int64, durationCount int64, durationMaxTryTime int64, durationMaxFailTimes int64, punishFactor ...float64) *RateLimiter {
	if durationUnitSec <= 0 || durationCount <= 0 || durationMaxTryTime <= 0 {
		panic("RateLimiter非法参数")
	}

	a := &RateLimiter{
		DuCount:             durationCount,
		DuUnitSec:           durationUnitSec,
		DuMaxTryTimes:       durationMaxTryTime,
		DuMaxFailTimes:      durationMaxFailTimes,
		PunishFactor:        0,
		TryTimesMap:         make(map[int64]int64, durationCount),
		ActiveDuUnitSec:     durationUnitSec,
		ActiveDuMaxTryTimes: durationMaxTryTime,
		LastPeriodTime:      gotime.UnixNowSec(),
		HasPunish:           false,
	}
	if len(punishFactor) > 0 {
		if punishFactor[0] < 0 || punishFactor[0] >= 1 {
			panic("RateLimiter 惩罚因子非法定义")
		}
		a.PunishFactor = punishFactor[0]
	}
	for i := int64(0); i < durationCount; i++ {
		a.TryTimesMap[i] = 0
	}
	return a
}

func (r *RateLimiter) rotate(duNum int64) {
	if duNum > 1 || r.TryTimesMap[0] < r.DuMaxFailTimes {
		// 重置
		r.ActiveDuUnitSec = r.DuUnitSec
		r.ActiveDuMaxTryTimes = r.DuMaxTryTimes
	}
	r.HasPunish = false
	//
	for duNum > 0 {
		for i := r.DuCount - 1; i > 0; i-- {
			r.TryTimesMap[i] = r.TryTimesMap[i-1]
		}
		r.TryTimesMap[0] = 0
		duNum--
	}
}

// 惩罚机制: 周期时间延时, 周期允许尝试次数缩小
func (r *RateLimiter) doPunish() {
	r.ActiveDuUnitSec += int64(float64(r.ActiveDuUnitSec) * r.PunishFactor)
	r.ActiveDuMaxTryTimes -= int64(float64(r.DuMaxTryTimes) * r.PunishFactor)
	if r.ActiveDuMaxTryTimes <= 0 {
		r.ActiveDuMaxTryTimes = 1
	}
	r.HasPunish = true
}

// 不通过返回等待恢复时间, 秒
func (r *RateLimiter) Pass() (pass bool, coolSec int64) {
	now := gotime.UnixNowSec()
	tsDiff := now - r.LastPeriodTime
	if tsDiff <= r.ActiveDuUnitSec {
		r.TryTimesMap[0] += 1
		t := r.TryTimesMap[0]
		if t <= r.DuMaxTryTimes {
			// do pass
			return true, 0
		}
		// do fail
		if t <= r.DuMaxFailTimes {
			coolSec = r.ActiveDuUnitSec - tsDiff
			return false, coolSec
		}
		// do fail and punish
		if r.DuMaxFailTimes <= 0 && !r.HasPunish {
			r.doPunish()
		}
		coolSec = r.ActiveDuUnitSec - tsDiff
		return false, coolSec
	}
	// do pass
	r.rotate(tsDiff / r.ActiveDuUnitSec)
	r.TryTimesMap[0] = 1
	r.LastPeriodTime = now
	return true, 0
}

// 取得周期热点列表
func (r *RateLimiter) GetTryMap() []int64 {
	var out = make([]int64, r.DuCount)
	for i, v := range r.TryTimesMap {
		out[i] = v
	}
	return out
}
