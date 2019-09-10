/*
 * Create by kidd
 * 2019-4-9
 */

package gotime

import (
	"time"
)

func UnixNow() int64 {
	return UnixNowSec()
}

func UnixNowSec() int64 {
	return time.Now().UnixNano() / time.Second.Nanoseconds()
}

func UnixNowMillSec() int64 {
	return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
}

func UnixNowNanoSec() int64 {
	return time.Now().UnixNano()
}

func ParseSeconds(sec int64) time.Time {
	return time.Unix(sec, 0)
}

func ParseMillSeconds(millSec int64) time.Time {
	return time.Unix(millSec/1000, millSec%time.Second.Nanoseconds())
}

func ParseNanoSeconds(nanoSec int64) time.Time {
	return time.Unix(nanoSec/time.Second.Nanoseconds(), nanoSec%time.Second.Nanoseconds())
}

func MillSecStrCST8(millSec int64, format ...string) string {
	var cstZone = time.FixedZone("CST", 8*3600) // 东8
	if format != nil && len(format) > 0 {
		return time.Unix(millSec/1000, 0).In(cstZone).Format(format[0])
	}
	return time.Unix(millSec/1000, 0).In(cstZone).Format("2006-01-02 15:04:05【CST8】")
}

func Time2StrCST8(t time.Time, format ...string) string {
	var cstZone = time.FixedZone("CST", 8*3600) // 东8
	if format != nil && len(format) > 0 {
		return t.In(cstZone).Format(format[0])
	}
	return t.In(cstZone).Format("2006-01-02 15:04:05【CST8】")
}

func TimeNow2StrCST8(format ...string) string {
	var cstZone = time.FixedZone("CST", 8*3600) // 东8
	if format != nil && len(format) > 0 {
		return time.Now().In(cstZone).Format(format[0])
	}
	return time.Now().In(cstZone).Format("2006-01-02 15:04:05【CST8】")
}

