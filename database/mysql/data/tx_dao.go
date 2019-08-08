// Mysql数据库操作增改查封装
//
// 所有方法外部通过data.操作数据库
//
// 由于golang不支持泛型，对于查询需要调用时配合
// Query及GetOne方法参数中的instance interface{}需要传入的是一个结构体类型，类似java的T.class。这里用：(*T)(nil)
// Query及GetOne返回的都是interface{}类型的指针，遍历时直接可以强转：t := it.Next().(*T)
// 具体可以看data_test.go测试Sample
//
// robot.guo

package data

import (
	"database/sql"
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/exwallet/go-common/database/mysql"
	"github.com/exwallet/go-common/gologger"
	"strings"
)

type TxDao struct {
	databaseKeys []string
	txs          map[string]*sql.Tx
}

//type MultiTxDao struct {
//	databases map[string]*TxDao
//}
//
//
//func NewMultiTxDao (txDaos... *TxDao) (nultiTxDao *MultiTxDao) {
//
//	return &MultiTxDao{
//		databases: nil,
//	}
//
//}

func NewTxDao(databaseKeys ...string) (dao *TxDao) {
	if len(databaseKeys) == 0 {
		return nil
	}
	var keys []string
	for _, v := range databaseKeys {
		keys = append(keys, v)
	}
	return &TxDao{
		databaseKeys: keys,
	}
}

// 调用在begin后
func (dao *TxDao) GetTx(databaseKey string) *sql.Tx {
	return dao.txs[databaseKey]
}

func (dao *TxDao) Query(tx *sql.Tx, queryPrepareSql string, instance interface{}, params ...interface{}) *arraylist.List {
	if tx == nil {
		return nil
	}
	stmt, err := tx.Prepare(queryPrepareSql)
	defer stmt.Close()
	if err != nil {
		gologger.Error("Query.db.Prepare()出错了：%v\n", err)
		return nil
	}
	rows, err := stmt.Query(params...)
	defer rows.Close()
	if err != nil {
		gologger.Error("Query.stmt.Query()出错了：%v\n", err)
		return nil
	}
	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		gologger.Error("Query.rows.Columns()出错了：%v\n", err)
	}
	//
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list := arraylist.New()

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			//panic(err.Error())
			continue
		}
		data := mapping(columns, values, instance)
		list.Add(data)
	}
	return list
}

func (dao *TxDao) GetOne(tx *sql.Tx, queryPrepareSql string, instance interface{}, params ...interface{}) interface{} {
	if tx == nil {
		return nil
	}
	if !strings.Contains(strings.ToLower(queryPrepareSql), "limit") {
		queryPrepareSql += " limit 0,1 "
	}
	list := dao.Query(tx, queryPrepareSql, instance, params...)
	if list.Size() == 1 {
		one, _ := list.Get(0)
		return one
		//return list.Get(0)
	}
	return nil
}

// 改了! 增加返回bool
func (dao *TxDao) Insert(tx *sql.Tx, oneSql *mysql.OneSql) (int64, bool) {
	if tx == nil {
		return -1, false
	}
	stmt, err := tx.Prepare(oneSql.Sql)
	defer stmt.Close()
	if err != nil {
		gologger.Error("Insert.db.Prepare()出错了：%v\n", err)
		return 0, false
	}
	res, err := stmt.Exec(oneSql.Params...)
	if err != nil {
		gologger.Error("Insert.stmt.Exec()出错了：%v\n", err)
		return 0, false
	}
	id, err := res.LastInsertId()
	if err != nil {
		gologger.Error("Insert.LastInsertId()出错了：%v\n", err)
		return 0, false
	}
	return id, true
}

// 改了! 增加返回bool
func (dao *TxDao) Update(tx *sql.Tx, updatePrepareSql string, params ...interface{}) (int64, bool) {
	if tx == nil {
		return -1, false
	}
	stmt, err := tx.Prepare(updatePrepareSql)
	defer stmt.Close()
	if err != nil {
		gologger.Error("Update.db.Prepare()出错了：%v\n", err)
		return 0, false
	}
	res, err := stmt.Exec(params...)
	if err != nil {
		gologger.Error("Update.stmt.Exec()出错了：%v\n", err)
		return 0, false
	}
	num, err := res.RowsAffected()
	if err != nil {
		gologger.Error("Update.RowsAffected()出错了：%v\n", err)
		return 0, false
	}
	if num <= 0 {
		return num, false
	}
	return num, true
}

func (dao *TxDao) Begin() error {

	if dao.databaseKeys == nil || len(dao.databaseKeys) == 0 {
		return fmt.Errorf("没有提供数据源,请初始化databaseKeys")
	}

	var dbs = map[string]*sql.DB{}
	var txs = map[string]*sql.Tx{}

	for _, v := range dao.databaseKeys {
		if dbs[v] != nil {
			continue
		}
		db := mysql.GetConnection(v)
		// 没有连接源
		if db == nil {
			return fmt.Errorf("databaseKey(%s)，找不到数据源，导致事务回滚\n", v)
		}
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("databaseKey(%s)，事务开启失败，导致事务回滚\n", v)
		}
		txs[v] = tx
	}

	// init
	dao.txs = txs

	return nil
}

func (dao *TxDao) Rollback(msg ...string) {
	s := ""
	if msg != nil || len(msg) > 0 {
		s = msg[0]
	}
	gologger.Error("sql发生回滚, 回滚信息:%s \n", s)
	rollback(dao.txs)
}

func (dao *TxDao) Commit() bool {
	for k, v := range dao.txs {
		err := v.Commit()
		if err != nil {
			gologger.Error("excute tx.Commit(): databaseKey(%s)，事物提交失败，导致事务回滚\n", k)
			return false
		}
	}
	return true
}
