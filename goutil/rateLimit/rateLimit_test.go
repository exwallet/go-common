/*
 * @Author: kidd
 * @Date: 8/10/19 6:27 PM
 */

package rateLimit

import (
	"github.com/exwallet/go-common/gologger"
	"testing"
	"time"
)

func Test_ratelimit(t *testing.T) {
	r := NewRateLimit(2, 5, 5, 8, 0.4)
	for {
		b, coolSec := r.Pass()
		gologger.Info("通过[%v], 冷静:%d 秒,  计数: %v", b, coolSec, r.GetTryMap())
		time.Sleep(time.Millisecond*100)
	}
}

