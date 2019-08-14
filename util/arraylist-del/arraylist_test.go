package arraylist_del

import (
	"github.com/exwallet/go-common/log"
	"testing"
)

type Stu struct {
	name string
}

type Person struct {
	name string
}

func Test_arraylist(t *testing.T) {

	list := new(ArrayList)

	list.Add("abc")
	list.Add("123")
	list.Add("3.14")
	list.Add(123)
	//
	stu1 := Stu{"a"}
	stu2 := Stu{"b"}
	stu3 := Stu{"c"}

	list.Add(stu1)
	list.Add(stu2)
	list.Add(stu3)
	//
	//fmt.Printf("%v\n", list)
	//
	//fmt.Println(list.Get(8))
	//
	//fmt.Println(list.Size())
	//
	//newList := list.Copy()
	//fmt.Println(newList)
	//
	//list.Set(2, 3.141526)
	//
	//fmt.Println(list)
	//
	//list.Insert(3, "333")
	//
	//fmt.Println(list)
	//
	//list.Remove(3)
	//
	//fmt.Println(list)
	//
	//list.RemoveByValue(3.141526)
	//
	//fmt.Println(list)
	//
	//list.RemoveByValue(stu2)
	//
	//fmt.Println(list)
	//
	//stu4 := list.GetPtr(3).(*Stu)
	//fmt.Printf("%T, %v\n", stu4, stu4)

	it := list.Iterator()
	for it.HasNext() {
		v := it.Next()
		//list.Set(it.cursor, v)
		log.Error("%T, %v, %p\n", v, v, &v)
	}
	//it = list.Iterator()
	//for it.HasNext(){
	//	v := it.Next()
	//	fmt.Printf("%T, %v, %p\n", v, v, &v)
	//}
}
