// list迭代器，配合list使用
// robot.guo

package arraylist_del

type Iterator struct {
	*ArrayList     // list引用地址
	cursor     int // 游标
}

// 是否存在元素
func (it *Iterator) HasNext() bool {
	if it.cursor < it.Size()-1 {
		return true
	}
	return false
}

// 获取下一个元素
func (it *Iterator) Next() interface{} {
	if it.cursor >= it.Size() {
		return nil
	}
	it.cursor++
	return it.Get(it.cursor)
}

// 获取当前游标位置
func (it *Iterator) GetCursor() int {
	return it.cursor
}

// 移除当前元素
func (it *Iterator) Remove() interface{} {
	if it.cursor == -1 {
		return nil
	}
	v := it.Get(it.cursor)
	it.ArrayList.Remove(it.cursor)
	it.cursor--
	return v
}
