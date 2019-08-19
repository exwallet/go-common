/*
 * @Author: kidd
 * @Date: 8/14/19 10:36 AM
 */

package limiter

import (
	"github.com/exwallet/go-common/goutil/gotime"
)

/**
限量器
*/

type NumberLimiter struct {
	DuUnitSec int64           `json:"duUnitSec"` // 周期时间,秒
	DuNum     int64           `json:"duNum"`     // 周期保存数量
	DuMaxNum  int64           `json:"duMaxNum"`  // 周期时间内最大数量
	NumMap    map[int64]int64 `json:"numMap"`    // 数量保存器
	LastTime  int64           `json:"lastTime"`  // 最后周期时间点
}

func NewNumberLimiter(duUnitSec int64, duNum int64, duMaxNum int64) *NumberLimiter {
	if duUnitSec <= 0 || duNum <= 0 || duMaxNum <= 0 {
		panic("NumberLimiter非法参数")
	}
	a := &NumberLimiter{
		DuUnitSec: duUnitSec,
		DuNum:     duNum,
		DuMaxNum:  duMaxNum,
		NumMap:    make(map[int64]int64, duNum),
		LastTime:  gotime.UnixNowSec(),
	}
	for i := int64(0); i < duNum; i++ {
		a.NumMap[i] = 0
	}
	return a
}

// 返回: 成功标志, 剩余数量额度, 到下次周期时间剩余时间,秒
func (l *NumberLimiter) Add(num int64) (succ bool, quotaLeft int64, secLeft int64) {
	now := gotime.UnixNowSec()
	tDiff := now - l.LastTime

	if tDiff <= l.DuUnitSec {
		// 未启用新周期
		secLeft = l.DuUnitSec - tDiff
		_n := l.DuMaxNum - l.NumMap[0]
		if num > _n {
			return false, _n, secLeft
		}
		l.NumMap[0] += num
		return true, l.DuMaxNum - l.NumMap[0], secLeft
	}
	// 新时间周期,数量超限
	if num > l.DuMaxNum {
		return false, l.DuMaxNum, l.DuUnitSec
	}

	// 新时间周期, 符合条件
	l.rotate(tDiff / l.DuUnitSec)
	l.LastTime = now
	l.NumMap[0] = num
	return true, l.DuMaxNum - num, l.DuUnitSec
}

// 撤消最后增加的数量
func (l *NumberLimiter) Revoke(num int64) {
	l.NumMap[0] -= num
}

func (l *NumberLimiter) rotate(duNum int64) {
	if duNum >= l.DuNum {
		for i := int64(0); i < duNum; i++ {
			l.NumMap[i] = 0
		}
		return
	}
	//
	for duNum > 0 {
		for i := l.DuNum - 1; i > 0; i-- {
			l.NumMap[i] = l.NumMap[i-1]
		}
		l.NumMap[0] = 0
		duNum--
	}
}

func (l *NumberLimiter) GetNumMap() []int64 {
	var out = make([]int64, l.DuNum)
	for i := int64(0); i < l.DuNum; i++ {
		out[i] = l.NumMap[i]
	}
	return out
}
