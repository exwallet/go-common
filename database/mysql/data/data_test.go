package data

import (
	"fmt"
	"github.com/exwallet/go-common/database/mysql"
	"testing"
)

type Per struct {
	Name string
}

type Stu struct {
	Id    int
	Name  string
	Age   int
	Addr  string
	Score float64
	per   Per
}

func (stu *Stu) service() bool {
	fmt.Println("============call me now================")
	return true
}

func Test_data(t *testing.T) {
	// configuration
	mysql.InitDataSources("../mysql.json")
	defer mysql.Close()

	// test insert
	//wg := sync.WaitGroup{}
	//for i := 0; i < 10000; i++ {
	//	wg.Add(1)
	//	go func(i int) {
	//		stu := Stu{0, "A" + strconv.Itoa(i), i, "Abc" + strconv.Itoa(i), 123.123, Per{"Per"}}
	//		id := Insert("main", stu)
	//		fmt.Println("id:", id)
	//		wg.Done()
	//	}(i)
	//}
	//wg.Wait()

	//test update
	//params := []interface{}{"newName222", 221, "newAddr222", 123.12345678 ,5}
	//c := Update("main", "update stu set name=?, age=?, addr=?, score=? where id=?", params...)
	//fmt.Println(c)

	// test query
	//list := Query("main", "select * from stu where id in(?,?,?,?,?)", (*Stu)(nil),1, 2, 3, 4, 5)
	//it := list.Iterator()
	//for it.HasNext() {
	//	v := it.Next().(*Stu)
	//	fmt.Printf("%p, %+v\n", v, v)
	//	//fmt.Println(v.Name, v.Age, v.Addr, v.Score)
	//}

	// test getOne
	//stu := GetOne("main", "select * from stu where id=?", (*Stu)(nil), 5).(*Stu)
	//fmt.Printf("%T, %+v\n", stu, stu)

	//test count
	//count := Count("main", "select * from stu where id in(?,?,?,?,?)", 1, 2, 3, 4, 5)
	//fmt.Println(count)

	// test transaction
	//sql1 := mysql.OneSql{"update stu set name=? where id=?", 1, []interface{}{"main1", 1}, "main"}
	//sql2 := mysql.OneSql{"update stu set name=? where id=?", 1, []interface{}{"main21", 25315}, "main2"}
	//list := arraylist.New()
	//list.Add(sql1, sql2)
	//// success
	//b := DoTrans(list)
	//fmt.Println(b)

	//sql1 := mysql.OneSql{"update stu set age=? where id=?", 1, []interface{}{12, 1}, "main"}
	//sql2 := mysql.OneSql{"update stu set age=? where id=?", 1, []interface{}{123, 25315}, "main2"}
	//list := arraylist.New()
	//list.Add(sql1, sql2)
	//// failed
	//b := DoTrans(list)
	//fmt.Println(b)

	//test transaction with other services
	//sql1 := mysql.OneSql{"update stu set age=? where id=?", 1, []interface{}{1234, 1}, "main"}
	//sql2 := mysql.OneSql{"update stu set age=? where id=?", 1, []interface{}{1234, 25315}, "main2"}
	//list :=  arraylist.New()
	//list.Add(sql1, sql2)
	//var stu = Stu{}
	//b := DoTransWithOtherService(list, &stu)
	//fmt.Println(b)
}
