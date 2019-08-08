/*
 * @Author: kidd
 * @Date: 1/19/19 10:24 PM
 */

package gois

import (
	"net"
	"regexp"
	"strings"
)

func IsInteger(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		str = strings.TrimSpace(str)
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		for _, v := range str {
			if v < '0' || v > '9' {
				return false
			}
		}
	}
	return true
}

func IsIp(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}
	return true
}

func IsEmail(s string) bool {
	pattern := `^[0-9A-Za-z][\.\-_0-9A-Za-z]*\@[0-9A-Za-z\-]+(\.[0-9A-Za-z]+)+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(s)
}

func IsCNMobile(s string) bool {
	reg := regexp.MustCompile(`^1[3456789]\d{9}$`)
	return reg.MatchString(s)
}
