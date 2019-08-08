/*
 * @Author: kidd
 * @Date: 1/19/19 10:24 PM
 */

package gosearch

import (
	"sort"
)

func SearchIntSlice(a []int, x int) int {
	sort.Ints(a)
	l := sort.SearchInts(a, x)
	if l == len(a) {
		return -1
	}
	return l
}

func SearchStringSlice(s []string, x string) int {
	sort.Strings(s)
	l := sort.SearchStrings(s, x)
	if l == len(s) {
		return -1
	}
	return l
}
