package redis

import (
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
	"testing"
)

type Stu struct {
	Name  string
	Age   int
	Score float64
	Sex   bool
}

func Test_redis(t *testing.T) {
	InitRedis("redis.json")

	//cluster, e := redis.NewCluster(&redis.Options{
	//	StartNodes:   []string{"192.168.3.111:7001", "192.168.3.112:7001", "192.168.3.113:7001"},
	//	ConnTimeout:  configuration.ConnTimeout,
	//	ReadTimeout:  configuration.ReadTimeout,
	//	WriteTimeout: configuration.WriteTimeout,
	//	KeepAlive:    1024,
	//	AliveTime:    30 * time.Second,
	//})
	//if e != nil {
	//	panic(e)
	//}
	//
	//fmt.Printf("%+v\n", cluster)

}

func Test_123(t *testing.T) {
	InitRedis("redis.json")

	fmt.Printf("Exists: %v\n", Exists("abc"))
	get, e := Get("abc")
	fmt.Printf("Get: %v %v\n", get, e)
	fmt.Printf("Set: %v\n", Set("abc", "abc"))
	fmt.Printf("SetAndExpire: %v\n", SetAndExpire("abc", "cbd", 30))
	fmt.Printf("Delete: %v\n", Delete("abc"))
	fmt.Printf("IncrBy: %v\n", IncrBy("ids"))

	var stu = Stu{"a", 18, 123.123, true}
	fmt.Printf("SetObj: %v\n", SetObj("abc", stu))
	fmt.Printf("SetObj: %v\n", SetObjAndExpire("abc", stu, 30))
	v, _ := GetObj("abc", (*Stu)(nil))
	fmt.Printf("%T, %+v\n", v, v)

}

func Test_1(t *testing.T) {
	InitRedis("redis.json")
	var stu2 = Stu{"b", 20, 123.343243, true}
	var stu3 = Stu{"c", 24, 12113.3432243, true}
	list1 := arraylist.New()
	list1.Add(stu2, stu3)
	fmt.Printf("setList: %v\n", SetList("abc", list1))
	fmt.Println(Get("abc"))
	l := GetList("abc", (*Stu)(nil))
	if l != nil {
		it := l.Iterator()
		for it.Next() {
			v := it.Value().(Stu)
			fmt.Printf("%T, %+v\n", v, v)
		}
		fmt.Printf("%T, %+v\n", l, l)
	}
	fmt.Printf("LPush: %v\n", LPush("xyz", -1, "a"))
	fmt.Printf("LPush: %v\n", LPush("xyz", -1, "b"))
	fmt.Printf("LPush: %v\n", LPush("xyz", -1, "c"))
	fmt.Printf("LPush: %v\n", LPush("xyz", -1, "c", "d", "e", "f", "g"))

	//l := LRange("xyz", 0, -1, nil)
	it := l.Iterator()
	for it.Next() {
		str := it.Value().(string)
		fmt.Println(str)
	}
}

func Test_2(t *testing.T) {
	InitRedis("redis.json")
	var stu2 = Stu{"b", 20, 123.343243, true}
	var stu3 = Stu{"c", 24, 12113.3432243, true}
	var stu4 = Stu{"d", 25, 113.343224343, true}
	var stu5 = Stu{"e", 26, 113.343122243, true}
	LPush("ttt", -1, stu2, stu3, stu4, stu5)
	l := LRange("ttt", 0, -1, (*Stu)(nil))
	it := l.Iterator()
	for it.Next() {
		stu := it.Value().(*Stu)
		fmt.Printf("%T, %+v\n", stu, stu)
	}

	fmt.Printf("LSET: %v\n", LSet("xyz", 3, "345", -1))
}

func Test_3(t *testing.T) {
	InitRedis("redis.json")
	var stu2 = Stu{"insert", 20, 123.343243, true}
	fmt.Printf("LSET: %v\n", LSet("xyz", 6, stu2, -1))

	fmt.Printf("RPop: %v\n", RPop("xyz", nil))
	stu := RPop("ttt", (*Stu)(nil))
	fmt.Printf("%T, %+v\n", stu, stu)

	fmt.Printf("LLength: %v\n", LLength("xyz"))

	fmt.Printf("HSet: %v\n", HSet("aaa", "a", "a", -1))
	fmt.Printf("HSet: %v\n", HSet("aaa1", "b", "b", -1))
	fmt.Printf("HSet: %v\n", HSet("aaa1", "c", "c", -1))
	fmt.Printf("HGet: %v\n", HGet("aaa1", "b", nil))
	fmt.Printf("HExists: %v\n", HExists("aaa1", "b"))
	fmt.Printf("HDelete: %v\n", HDelete("aaa1", "b"))

}

func Test_4(t *testing.T) {
	var stu2 = Stu{"stu2", 20, 123.343243, true}
	var stu3 = Stu{"stu3", 24, 12113.3432243, false}
	var stu4 = Stu{"stu4", 25, 113.343224343, false}
	fmt.Printf("HSet: %v\n", HSet("xxx", "stu2", stu2, -1))
	fmt.Printf("HSet: %v\n", HSet("xxx", "stu3", stu3, -1))
	fmt.Printf("HSet: %v\n", HSet("xxx", "stu4", stu4, -1))
	stu := HGet("xxx", "stu3", (*Stu)(nil))
	fmt.Printf("%T, %+v\n", stu, stu)

	v := HGetAll("xxx", (*Stu)(nil))
	fmt.Printf("%T, %v", v, v)
	for k, val := range v {
		fmt.Printf("%s, %T, %+v\n", k, val, val)
		stu := val.(*Stu)
		fmt.Println(stu.Name)
	}
}
