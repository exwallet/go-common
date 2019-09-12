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
	DuUnitSec           int64           `json:"duUnitSec"`           // 周期时长
	DuMaxTryTimes       int64           `json:"duMaxTryTimes"`       // 周期内允许尝试次数
	DuMaxFailTimes      int64           `json:"duMaxFailTimes"`      // 周期内允许最大失败次数
	CalmIntervalSec     int64           `json:"calmIntervalSec"`     // 每次调用冷静时间, <=0 不需要冷静
	PunishFactor        float64         `json:"punishFactor"`        // (0<n<1)惩罚倍数因子,失败次数达多少时, 周期变长, 允许次数变小 .
	TryTimes            int64           `json:"tryTimes"`            //
	ActiveDuUnitSec     int64           `json:"activeDuUnitSec"`     //
	ActiveDuMaxTryTimes int64           `json:"activeDuMaxTryTimes"` //
	LastPeriodTime      int64           `json:"lastPeriodTime"`      // 最新周期时间点
	HasPunish           bool            `json:"hasPunish"`           //
}

// durationMaxFailTimes <= 0 不惩罚
func NewRateLimiter(durationUnitSec int64, durationMaxTryTime int64, durationMaxFailTimes int64, calmIntervalSec int64, punishFactor ...float64) *RateLimiter {
	if durationUnitSec <= 0 || durationMaxTryTime <= 0 {
		panic("RateLimiter非法参数")
	}

	a := &RateLimiter{
		DuUnitSec:           durationUnitSec,
		DuMaxTryTimes:       durationMaxTryTime,
		DuMaxFailTimes:      durationMaxFailTimes,
		CalmIntervalSec:     calmIntervalSec,
		PunishFactor:        0,
		TryTimes:            0,
		ActiveDuUnitSec:     durationUnitSec,
		ActiveDuMaxTryTimes: durationMaxTryTime,
		LastPeriodTime:      gotime.UnixNowSec(),
		HasPunish:           false,
	}
	if len(punishFactor) > 0 {
		a.PunishFactor = punishFactor[0]
	}
	return a
}

func (r *RateLimiter) rotate() {
	if r.TryTimes < r.DuMaxTryTimes + r.DuMaxFailTimes {
		r.ActiveDuUnitSec = r.DuUnitSec
		r.ActiveDuMaxTryTimes = r.DuMaxTryTimes
	}
	r.HasPunish = false
}

// 惩罚机制: 周期时间延时, 周期允许尝试次数缩小
func (r *RateLimiter) doPunish() {
	if r.PunishFactor <= 0 {
		return
	}
	r.ActiveDuUnitSec += int64(float64(r.ActiveDuUnitSec) * r.PunishFactor)
	r.ActiveDuMaxTryTimes -= int64(float64(r.DuMaxTryTimes) * r.PunishFactor)
	if r.ActiveDuMaxTryTimes <= 0 {
		r.ActiveDuMaxTryTimes = 1
	}
	r.HasPunish = true
}

// 不通过返回等待恢复时间, 秒
func (r *RateLimiter) Do() (pass bool, coolSec int64) {
	now := gotime.UnixNowSec()
	tsDiff := now - r.LastPeriodTime
	if tsDiff <= r.ActiveDuUnitSec {
		r.TryTimes += 1
		if r.TryTimes <= r.DuMaxTryTimes {
			// do pass, but check calm
			if r.CalmIntervalSec > 0 && tsDiff < r.CalmIntervalSec {
				return false, r.CalmIntervalSec - tsDiff
			}
			return true, 0
		}
		// do fail
		if r.TryTimes <= r.DuMaxFailTimes {
			coolSec = r.ActiveDuUnitSec - tsDiff
			return false, coolSec
		} else if !r.HasPunish {
			r.doPunish()
		}
		return false, r.ActiveDuUnitSec - tsDiff
	}
	// do pass
	r.rotate()
	r.TryTimes = 1
	r.LastPeriodTime = now
	return true, 0
}

func (r *RateLimiter) Revoke() {
	r.TryTimes -= 1
}

