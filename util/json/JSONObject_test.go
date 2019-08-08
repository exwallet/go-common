package json

import (
	"fmt"
	"testing"
)

type Stu struct {
	Name string
}

func Test_JSONObject(t *testing.T) {

	jo := NewJSONObject()
	jo.Put("a", "a")
	jo.Put("b", "b")
	jo.Put("c", "c")
	jo.Put("e", "e")
	jo.Put("f", "f")
	jo.Put("g", "g")
	jo.Put("h", "h")
	jo.Put("v", 1)
	jo.Put("x", 1)
	jo.Put("h", 2)

	fmt.Printf("%T, %v\n", jo, jo)

	fmt.Println(jo.Get("a"))

	jo.Remove("b")

	vv := jo.Get("a")
	fmt.Println(jo)
	vv = vv.(string)
	fmt.Println(jo)

	fmt.Printf("%T, %v\n", jo, jo)

	for k, v := range jo.Elements() {
		fmt.Println(k, "-->", v)
	}

	fmt.Println(jo.ToJSONString())

	stu1 := Stu{"a"}
	stu2 := Stu{"b"}
	stu3 := Stu{"c"}

	jo1 := NewJSONObject()
	jo1.Put("a", stu1)
	jo1.Put("b", stu2)
	jo1.Put("c", stu3)

	fmt.Println(jo1)

	fmt.Println(jo1.ToJSONString())

	str := "{\"name\" : \"hahaha\"}"

	jo, err := ParseJSONObject(str, (*Stu)(nil))
	fmt.Println(err)
	fmt.Printf("%T, %+v\n", jo, jo)
	fmt.Println(jo.ToJSONString())
}
