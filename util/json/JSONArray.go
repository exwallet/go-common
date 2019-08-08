//
// 根据arraylist封装成类似于java JSONArray工具包
// JSONArray元素只支持JSONObject
//
// 遍历map，无序
// for i := 0; i < array.Size(); i++ {
//		jo := array.Get(i)
// }
//
// robot.guo

package json

import (
	"bytes"
	"encoding/json"
	"github.com/exwallet/go-common/gods/lists/arraylist"
	"github.com/exwallet/go-common/util/strutil"
	"strings"
)

type JSONArray struct {
	datas *arraylist.List
}

func NewJSONArray() *JSONArray {
	ja := new(JSONArray)
	ja.datas = arraylist.New()
	return ja
}

func ParseJSONArray(str string) *JSONArray {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	//str = strings.TrimSpace(str)
	if strings.Index(str, "[") != 0 {
		panic("array must be start with [")
	}
	if strings.Index(str, "]") != len(str)-1 {
		panic("array must be end with ]")
	}
	str = strutil.Substring(str, 1, len(str)-1)
	str = strings.Replace(str, "},{", "}`{", -1)
	arr := strings.Split(str, "`")
	list := arraylist.New()
	for _, v := range arr {
		var mapResult map[string]interface{}
		if err := json.Unmarshal([]byte(v), &mapResult); err != nil {
			panic(err)
		}
		jo := NewJSONObjectWithDatas(mapResult)
		list.Add(jo)
	}
	ja := NewJSONArray()
	ja.datas = list
	return ja
}

func ToJSONArray(datas *arraylist.List, instance interface{}) *JSONArray {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	ja := NewJSONArray()
	it := datas.Iterator()
	for it.Next() {
		data := it.Value()
		jo, err := ToJSONObject(data)
		if err != nil {
			panic(err)
		}
		ja.Add(jo)
	}
	return ja
}

func (ja *JSONArray) Size() int {
	return ja.datas.Size()
}

func (ja *JSONArray) Add(jsonObject *JSONObject) {
	ja.datas.Add(jsonObject)
}

func (ja *JSONArray) Get(index int) *JSONObject {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	if o, b := ja.datas.Get(index); b {
		return o.(*JSONObject)
	}
	return nil
}

func (ja *JSONArray) Remove(index int) {
	ja.datas.Remove(index)
}

func (ja *JSONArray) ToJSONString() string {
	var results bytes.Buffer
	results.WriteString("[")
	it := ja.datas.Iterator()
	for it.Next() {
		jo := it.Value().(*JSONObject)
		results.WriteString(jo.ToJSONString())
		if it.Index()+1 == ja.datas.Size() {
			results.WriteString(",")
		}
	}
	results.WriteString("]")
	return results.String()
}
