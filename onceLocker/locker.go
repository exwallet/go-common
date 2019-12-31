/*
 * @Author: kidd
 * @Date: 2019/12/31 下午5:46
 */

package onceLocker

import (
	"github.com/exwallet/go-common/log"
	"sync"
)

// 占用锁
type OnceLocker struct {
	isLock bool
	*sync.Mutex
}

func NewOnceLocker() *OnceLocker {
	return &OnceLocker{
		isLock: false,
		Mutex:  new(sync.Mutex),
	}
}

func (l *OnceLocker) IsLock() bool {
	return l.isLock
}

func (l *OnceLocker) Lock() {
	if l.isLock {
		log.Error("already locked")
		return
	}
	l.Mutex.Lock()
	l.isLock = true
}

func (l *OnceLocker) Unlock() {
	if l.isLock {
		l.Mutex.Unlock()
		l.isLock = false
		return
	}
	log.Error("not lock yet")
}


