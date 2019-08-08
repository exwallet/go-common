//
// 根据golang的map封装成类似于java JSONObject工具包
//
// 遍历map，无序
// for i, v := range jo.Elements {
//
// }
//
// robot.guo

package json

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type JSONObject struct {
	datas map[string]interface{}
}

func NewJSONObject() *JSONObject {
	jo := new(JSONObject)
	jo.datas = make(map[string]interface{})
	return jo
}

func NewJSONObjectWithDatas(datas map[string]interface{}) *JSONObject {
	jo := new(JSONObject)
	jo.datas = datas
	return jo
}

func ParseJSONObject(str string, instance interface{}) (*JSONObject, error) {
	ptr := reflect.New(reflect.TypeOf(instance).Elem()).Interface()
	err := json.Unmarshal([]byte(str), ptr)
	if err != nil {
		return nil, err
	}
	return ToJSONObject(ptr)
}

func ToJSONObject(data interface{}) (*JSONObject, error) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var dataMap = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		dataMap[t.Field(i).Name] = v.Field(i).Interface()
	}
	jo := NewJSONObject()
	jo.datas = dataMap
	return jo, nil
}

func (jo *JSONObject) Size() int {
	return len(jo.datas)
}

func (jo *JSONObject) Elements() map[string]interface{} {
	return jo.datas
}

func (jo *JSONObject) Put(k string, v interface{}) {
	jo.datas[k] = v
}

func (jo *JSONObject) Get(k string) interface{} {
	if v, ok := jo.datas[k]; ok {
		return v
	}
	return nil
}

func (jo *JSONObject) Remove(k string) {
	delete(jo.datas, k)
}

func (jo *JSONObject) ContainsKey(k string) bool {
	if _, ok := jo.datas[k]; ok {
		return true
	}
	return false
}

func (jo *JSONObject) ToJSONString() string {
	str, err := json.Marshal(jo.datas)
	if err != nil {
		fmt.Println("JSONObject toJSONString err : ", err)
		return "{}"
	}
	return string(str)
}
