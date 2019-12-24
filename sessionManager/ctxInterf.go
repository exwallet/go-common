/*
 * @Author: kidd
 * @Date: 12/24/19 3:35 PM
 */

package sessionManager


type ContextInf interface {
	GetCookie(key string) string
	SetCookie(key string, val string, ageSeconds int64, path string, domain string, secure bool, HttpOnly bool)
	IP() string
	GetHeader(key string) string
}

