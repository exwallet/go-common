/*
 * @Author: kidd
 * @Date: 2019/12/31 下午5:47
 */

package onceLocker

import (
	"fmt"
	"testing"
)



func Test_1(t *testing.T) {
	l1 := NewOnceLocker()
	fmt.Println("l1 before lock", l1.isLock)
	l1.Lock()
	fmt.Println("l1 after lock ", l1.isLock)
	l1.Lock()
	l1.Unlock()
	fmt.Println("l1 after unlock ", l1.isLock)

	l2 := NewOnceLocker()
	fmt.Println("l2 before lock", l2.isLock)
	l2.Lock()
	fmt.Println("l2 after lock ", l2.isLock)
	l2.Lock()
	l2.Unlock()
	fmt.Println("l2 after unlock ", l2.isLock)




}

