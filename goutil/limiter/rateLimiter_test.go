/*
 * @Author: kidd
 * @Date: 8/10/19 6:27 PM
 */

package limiter

import (
	"github.com/exwallet/go-common/log"
	"testing"
	"time"
)

func Test_ratelimit(t *testing.T) {
	r := NewRateLimiter(2, 5, 5, 8, 0.4)
	for {
		b, coolSec := r.Pass()
		log.Info("通过[%v], 冷静:%d 秒,  计数: %v", b, coolSec, r.GetTryMap())
		time.Sleep(time.Millisecond*100)
	}
}

func Test_numberLimiter(t *testing.T) {
	l := NewNumberLimiter(5, 10, 20)
	for {
		succ, quotaLeft, timeLeft := l.Add(1)
		log.Info("通过[%v], 剩余额度[%v], 剩余时间[%v]秒, 数量表[%v]", succ, quotaLeft, timeLeft, l.GetNumMap())
		time.Sleep(time.Millisecond*200)
	}
}
