//
// 根据golang的slice封装成类似于java ArrayList的用法
//
// 普通遍历list
// for i := 0; i < list.Size(); i++{
//
// }
//
// 迭代器便利list
// it := l.Iterator()
// for it.HasNext() {
//	 v := it.Next()
//   it.Remove()	// remove不会越界
// }
//
// arraylist返回的interface{}都是一份拷贝，如果需要影响arraylist，需要再次调用Set(index, value)
// robot.guo

package arraylist_del

import "sync"

type ArrayList struct {
	elements []interface{}
	lock     sync.Mutex
}

// create a new ArrayList
func New() *ArrayList {
	return new(ArrayList)
}

// create a new ArrayList
func NewWithElements(elements []interface{}) *ArrayList {
	l := new(ArrayList)
	l.elements = elements
	return l
}

// create a new Iterator
func (l *ArrayList) Iterator() *Iterator {
	return &Iterator{l, -1}
}

// 获取List的长度
func (l *ArrayList) Size() int {
	return len(l.elements)
}

// 清空List
func (l *ArrayList) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.elements = make([]interface{}, 0)
}

// 拷贝一份新的List，深拷贝
func (l *ArrayList) Copy() *ArrayList {

	l.lock.Lock()
	defer l.lock.Unlock()

	newList := New()
	els := make([]interface{}, len(l.elements), cap(l.elements))
	copy(els, l.elements)
	newList.elements = els

	return newList
}

//
func (l *ArrayList) GetElements() []interface{} {
	return l.elements
}

// 获取元素
func (l *ArrayList) Get(index int) interface{} {

	l.lock.Lock()
	defer l.lock.Unlock()

	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	//if index < 0 || index >= l.Size() {
	//	return nil
	//}

	return l.elements[index]
}

// 更新元素
func (l *ArrayList) Set(index int, e interface{}) bool {

	l.lock.Lock()
	defer l.lock.Unlock()

	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	if e == nil {
		return false
	}

	//
	//if index < 0 || index >= l.Size() {
	//	return false
	//}

	l.elements[index] = e

	return true
}

// 添加元素
func (l *ArrayList) Add(e ...interface{}) bool {

	l.lock.Lock()
	defer l.lock.Unlock()

	if e == nil {
		return false
	}

	l.elements = append(l.elements, e...)

	return true
}

// 添加一个List
func (l *ArrayList) AddAll(list *ArrayList) bool {

	l.lock.Lock()
	defer l.lock.Unlock()

	if list == nil {
		return false
	}

	l.elements = append(l.elements, (*list).elements...)

	return true
}

// 指定位置插入
func (l *ArrayList) Insert(index int, e interface{}) bool {

	l.lock.Lock()
	defer l.lock.Unlock()

	defer func() bool {
		recover()
		return false
	}()

	if e == nil {
		return false
	}

	size := l.Size()

	switch {
	case size == 0 && index == 0:
		l.elements = append(l.elements, e)
	case size == index+1:
		l.elements = append(l.elements, e)
	case index == 0:
		l.elements = append([]interface{}{e}, l.elements...)
	case size < index:
		es := make([]interface{}, index+1, index*2)
		copy(es, l.elements)
		es[index] = e
		l.elements = es
	default:
		es := make([]interface{}, 0, size*2)
		es = append(es, l.elements[:index]...)
		es = append(es, e)
		es = append(es, l.elements[index:]...)
		l.elements = es
	}

	return true
}

// 移除元素
func (l *ArrayList) Remove(index int) bool {

	l.lock.Lock()
	defer l.lock.Unlock()

	if index < 0 || index >= l.Size() {
		return false
	}

	size := l.Size()

	switch {
	case index == 0:
		l.elements = l.elements[1:]
	case index == size-1:
		l.elements = l.elements[:size-1]
	default:
		l.elements = append(l.elements[:index], l.elements[index+1:]...)
	}
	return true
}

// 移除元素
func (l *ArrayList) RemoveByValue(v interface{}) bool {
	// 迭代器
	it := l.Iterator()
	for it.HasNext() {
		if it.Next() == v {
			it.Remove()
			return true
		}
	}
	return false
}
