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
	r := NewRateLimit(5, 8, 2, 4, 0.5)
	for {
		b, coolSec := r.Pass()
		v := 0
		if b {
			v = 1
		}
		gologger.Info("%v, 冷静:%d 秒,  命中表: %v", v, coolSec, r.GetTryMap())
		time.Sleep(time.Second)
	}

}

