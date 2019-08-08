/*
 * @Author: kidd
 * @Date: 1/19/19 10:24 PM
 */

package goslice

func InSlice(v interface{}, sl []interface{}) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InSliceInt(v int, sl []int) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InSliceString(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func SliceIntersect(slice1, slice2 []interface{}) (diffslice []interface{}) {
	for _, v := range slice1 {
		if InSlice(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

func SliceIntersectInt(slice1, slice2 []int) (diffslice []int) {
	for _, v := range slice1 {
		if InSliceInt(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}
