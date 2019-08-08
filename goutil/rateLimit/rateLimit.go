/*
 * @Author: kidd
 * @Date: 7/30/19 10:52 AM
 */

package rateLimit

import (
	"sync"
	"time"
)

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
