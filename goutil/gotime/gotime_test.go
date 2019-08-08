/*
 * @Author: kidd
 * @Date: 7/17/19 2:26 PM
 */

package gotime

import (
	"github.com/exwallet/go-common/gologger"
	"testing"
)

func Test_MillSecStrCST8(t *testing.T){
	now := UnixNowMillSec()
	s := MillSecStrCST8(now)
	gologger.Info(s)
}
