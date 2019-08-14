/*
 * @Author: kidd
 * @Date: 8/14/19 10:36 AM
 */

package limiter

import (
	"github.com/exwallet/go-common/goutil/gotime"
	"sync"
)

/**
限量器
*/

type NumberLimiter struct {
	DuUnitSec int64           // 周期时间,秒
	DuNum     int64           // 周期保存数量
	DuMaxNum  int64           // 周期时间内最大数量
	numMap    map[int64]int64 // 数量保存器
	lastTime  int64           // 最后周期时间点
	lock      sync.Mutex
}

func NewNumberLimiter(duUnitSec int64, duNum int64, duMaxNum int64) *NumberLimiter {
	if duUnitSec <=0 || duNum <= 0|| duMaxNum <= 0 {
		panic("NumberLimiter非法参数")
	}
	a := &NumberLimiter{
		DuUnitSec: duUnitSec,
		DuNum:     duNum,
		DuMaxNum:  duMaxNum,
		numMap:    make(map[int64]int64, duNum),
		lastTime:  gotime.UnixNowSec(),
		lock:      sync.Mutex{},
	}
	for i := int64(0); i < duNum; i++ {
		a.numMap[i] = 0
	}
	return a
}

// 返回: 成功标志, 剩余数量额度, 到下次周期时间剩余时间,秒
func (l *NumberLimiter) Add(num int64) (succ bool, quotaLeft int64, secLeft int64) {
	now := gotime.UnixNowSec()
	tDiff := now - l.lastTime

	if tDiff <= l.DuUnitSec {
		// 未启用新周期
		secLeft = l.DuUnitSec - tDiff
		_n := l.DuMaxNum - l.numMap[0]
		if num > _n {
			return false, _n, secLeft
		}
		l.numMap[0] += num
		return true, l.DuMaxNum - l.numMap[0], secLeft
	}
	// 新时间周期,数量超限
	if num > l.DuMaxNum {
		return false, l.DuMaxNum, l.DuUnitSec
	}
	// 新时间周期, 符合条件
	l.rotate()
	l.lastTime = now
	l.numMap[0] = num
	return true, l.DuMaxNum - num, l.DuUnitSec
}

func (l *NumberLimiter) rotate() {
	for i := l.DuNum - 1; i > 0; i-- {
		l.numMap[i] = l.numMap[i-1]
	}
	l.numMap[0] = 0
}

func (l *NumberLimiter) GetNumMap() []int64 {
	var out = make([]int64, l.DuNum)
	for i := int64(0); i < l.DuNum; i++ {
		out[i] = l.numMap[i]
	}
	return out
}
