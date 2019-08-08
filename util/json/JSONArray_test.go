package json

import (
	"fmt"
	"github.com/exwallet/go-common/gods/lists/arraylist"
	"testing"
)

type Stu2 struct {
	Name string
}

func Test_JSONArray(t *testing.T) {

	jo := NewJSONObject()
	jo.Put("a", "a")
	jo.Put("b", "b")
	jo.Put("c", "c")
	jo.Put("e", "e")

	jo2 := NewJSONObject()
	jo2.Put("h", "h")
	jo2.Put("j", "j")
	jo2.Put("k", "k")
	jo2.Put("l", "l")

	ja := NewJSONArray()
	ja.Add(jo)
	ja.Add(jo2)

	fmt.Println(ja.ToJSONString())
	fmt.Println(ja.Size())

	//jo3 := ja.Get(1)
	//fmt.Println(jo3.ToJSONString())

	//ja.Remove(0)
	fmt.Println(ja.ToJSONString())

	stu1 := Stu2{"a"}
	stu2 := Stu2{"b"}
	stu3 := Stu2{"c"}

	stu5 := Stu2{"d"}
	stu4 := Stu2{"e"}
	stu6 := Stu2{"f"}

	jo3 := NewJSONObject()
	jo3.Put("a", stu1)
	jo3.Put("b", stu2)
	jo3.Put("c", stu3)

	jo4 := NewJSONObject()
	jo4.Put("d", stu5)
	jo4.Put("e", stu4)
	jo4.Put("f", stu6)

	ja2 := NewJSONArray()
	ja2.Add(jo3)
	ja2.Add(jo4)

	fmt.Println(ja2.ToJSONString())
	ja2.Remove(1)
	fmt.Println(ja2.Size())
	fmt.Println(ja2.ToJSONString())

	str := "[{\"a\":{\"Name\":\"a\"},\"b\":{\"Name\":\"b\"},\"c\":{\"Name\":\"c\"}},{\"d\":{\"Name\":\"d\"},\"e\":{\"Name\":\"e\"},\"f\":{\"Name\":\"f\"}}]"

	ja3 := ParseJSONArray(str)
	fmt.Println("datas->", ja3.ToJSONString())

	list := arraylist.New()
	list.Add(stu1, stu2, stu3, stu4, stu5)
	ja4 := ToJSONArray(list, (*Stu2)(nil))
	fmt.Println(ja4.ToJSONString())
}
