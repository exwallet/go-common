/*
 * Create by kidd
 * 2019-2-24
 */

package data

import (
	"database/sql"
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/exwallet/go-common/database/mysql"
	"github.com/exwallet/go-common/log"
	"github.com/exwallet/go-common/util/strutil"
	"strings"
)

//const PrimaryKeyName = "primary"

type Model struct {
}

// ******************************  包装开始 ********************************  改为用 gods的arraylist结构
type Dao struct {
	DatabaseKey string
}

func NewDao(databaseKey string) *Dao {
	dao := &Dao{
		DatabaseKey: databaseKey,
	}
	return dao
}

// new
//
func (d *Dao) Query1(result interface{}, queryPrepareSql string, instance interface{}, params ...interface{}) {

}

func (d *Dao) Query(queryPrepareSql string, instance interface{}, params ...interface{}) (query *arraylist.List, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Query[%s] 捕获到panic：%v", queryPrepareSql, e)
		}
	}()
	//return Query(d.DatabaseKey, queryPrepareSql, instance, params...)
	db := mysql.GetConnection(d.DatabaseKey)
	if db == nil {
		return nil, fmt.Errorf("cannot found out datasource, database(%s)\n", d.DatabaseKey)
	}
	stmt, e := db.Prepare(queryPrepareSql)
	defer stmt.Close()
	if e != nil {
		return nil, fmt.Errorf("Query[%s] 出错了：%v\n", queryPrepareSql, e)
	}
	rows, e := stmt.Query(params...)
	defer rows.Close()
	if e != nil {
		return nil, fmt.Errorf("Query[%s] 出错了：%v\n", queryPrepareSql, e)
	}
	// 获取列名
	columns, e := rows.Columns()
	if e != nil {
		return nil, fmt.Errorf("Query.rows.Columns()出错了：%v\n", e)
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
			continue
		}
		data := mapping(columns, values, instance)
		list.Add(data)
	}
	return list, nil
}
func (d *Dao) QueryByPage(queryPrepareSql string, pageNo int, pageSize int, instance interface{}, params ...interface{}) (query *arraylist.List, err error) {
	//logger.Debug("sql : %s", queryPrepareSql)
	//logger.Debug("sql params:   %d, %+v", len(params), params)
	if strings.Index(queryPrepareSql, "limit") < 0 {
		if pageNo < 1 {
			pageNo = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		limit := fmt.Sprintf(" limit %d,%d ", pageSize*(pageNo-1), pageSize)
		queryPrepareSql += limit
	}
	return d.Query(queryPrepareSql, instance, params...)
}

func (d *Dao) GetOne(queryPrepareSql string, instance interface{}, params ...interface{}) (one interface{}, err error) {
	//return GetOne(d.DatabaseKey, queryPrepareSql, instance, params...)
	db := mysql.GetConnection(d.DatabaseKey)
	if db == nil {
		return nil, fmt.Errorf("cannot found out datasource, database(%s)", d.DatabaseKey)
	}
	if !strings.Contains(strings.ToLower(queryPrepareSql), "limit") {
		queryPrepareSql += " limit 0,1 "
	}
	list, err := d.Query(queryPrepareSql, instance, params...)
	if err != nil {
		return
	}
	if list == nil {
		return nil, nil
	}
	if list.Size() == 1 {
		one, _ := list.Get(0)
		return one, nil
	}
	return nil, nil
}

func (d *Dao) Count(queryPrepareSql string, params ...interface{}) int64 {
	//return Count(d.DatabaseKey, queryPrepareSql, params...)
	db := mysql.GetConnection(d.DatabaseKey)
	if db == nil {
		log.Error("cannot found out datasource, database(%s)", d.DatabaseKey)
		return -1
	}
	if strings.Contains(strings.ToLower(queryPrepareSql), "group by") {
		queryPrepareSql = "select count(*) from (" + queryPrepareSql + ") t"
	} else {
		queryPrepareSql = "select count(*) " + strutil.Substring(queryPrepareSql, strings.Index(strings.ToLower(queryPrepareSql), "from"), len(queryPrepareSql))
	}
	var count int64
	row := db.QueryRow(queryPrepareSql, params...)
	if row == nil {
		log.Error("Count.db.QueryRow()出错了()")
		return 0
	}
	row.Scan(&count)
	return count

}

func (d *Dao) Insert(v interface{}) (id int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("捕获到panic：%v", e)
		}
	}()
	//return Insert(d.DatabaseKey, v)
	db := mysql.GetConnection(d.DatabaseKey)
	if db == nil {
		return -1, fmt.Errorf("cannot found out datasource, database(%s)", d.DatabaseKey)
	}
	oneSql, table, err := makeInsertSql(v, d.DatabaseKey)
	if err != nil {
		return 0, fmt.Errorf("[%s]生成Sql出错:%s", table, err.Error())
	}
	stmt, err := db.Prepare(oneSql.Sql)
	defer stmt.Close()
	if err != nil {
		return 0, fmt.Errorf("[%s]Insert.db.Prepare出错了(%s):  %v", table, oneSql.Sql, err)
	}
	res, err := stmt.Exec(oneSql.Params...)
	if err != nil {
		return 0, fmt.Errorf("[%s]Insert.stmt.Exec出错了()：%v", table, err)
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("[%s]Insert.LastInsertId出错了()：%v", table, err)
	}
	return id, nil
}

func (d *Dao) InsertByOneSql(oneSql *mysql.OneSql) (id int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("捕获到panic：%v", e)
		}
	}()
	db := mysql.GetConnection(oneSql.DatabaseKey)
	if db == nil {
		return -1, fmt.Errorf("cannot found out datasource, database(%s)", oneSql.DatabaseKey)
	}
	stmt, err := db.Prepare(oneSql.Sql)
	defer stmt.Close()
	if err != nil {
		return 0, fmt.Errorf("InsertByOneSql.db.Prepare()出错了：%v\n", err)
	}
	res, err := stmt.Exec(oneSql.Params...)
	if err != nil {
		return 0, fmt.Errorf("InsertByOneSql.stmt.Exec()出错了：%v\n", err)
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertByOneSql.LastInsertId()出错了：%v\n", err)
	}
	return id, nil
}

func (d *Dao) Update(updatePrepareSql string, params ...interface{}) (num int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("捕获到panic：%v", e)
		}
	}()
	//return Update(d.DatabaseKey, updatePrepareSql, params...)
	db := mysql.GetConnection(d.DatabaseKey)
	//gologger.Debug("sql: %s", updatePrepareSql)
	//gologger.Debug("params: %+v", params)
	if db == nil {
		return -1, fmt.Errorf("cannot found out datasource, database(%s)\n", d.DatabaseKey)
	}
	stmt, err := db.Prepare(updatePrepareSql)
	defer stmt.Close()
	if err != nil {
		return 0, fmt.Errorf("Update.db.Prepare()出错了：%v\n", err)
	}
	res, err := stmt.Exec(params...)
	if err != nil {
		return 0, fmt.Errorf("Update.stmt.Exec()出错了：%v\n", err)
	}
	num, err = res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Update.RowsAffected()出错了：%v\n", err)
	}
	if num <= 0 {
		return num, fmt.Errorf("Update.更新记录数为0")
	}
	return num, nil
}

func (d *Dao) DoTrans(sqls []*mysql.OneSql) (success bool) {
	// sqls 不能有nil
	for _, s := range sqls {
		if s == nil {
			log.Error("事务中包含nil, 导致事务回滚")
			return false
		}
	}
	//return DoTrans(sqls)
	// var
	var txs map[string]*sql.Tx = map[string]*sql.Tx{}
	// execute
	success = execute(txs, sqls)
	// 执行失败了
	if !success {
		return
	}
	// commit
	for k, v := range txs {
		err := v.Commit()
		if err != nil {
			log.Error("excute tx.Commit(): databaseKey(%s)，事物提交失败，导致事务回滚\n", k)
			rollback(txs)
			return
		}
	}
	// execute success
	success = true
	return
}

func (d *Dao) DoTransWithOtherService(sqls []*mysql.OneSql, callWithMysqlTx CallWithMysqlTx) (success bool) {
	//return DoTransWithOtherService(sqls, callWithMysqlTx)
	// var
	var txs map[string]*sql.Tx = map[string]*sql.Tx{}
	// execute
	success = execute(txs, sqls)
	// 执行失败了
	if !success {
		return
	}
	// other service
	if !callWithMysqlTx.service() {
		success = false
		rollback(txs)
		return
	}
	// commit
	for k, v := range txs {
		err := v.Commit()
		if err != nil {
			log.Error("excute tx.Commit(): databaseKey(%s)，事物提交失败，导致事务回滚\n", k)
			rollback(txs)
			return
		}
	}
	// execute success
	success = true
	return
}

func (d *Dao) NewOneSql(sql string, efRows int, params ...interface{}) *mysql.OneSql {
	return &mysql.OneSql{
		Sql:         sql,
		EfRows:      efRows,
		Params:      params,
		DatabaseKey: d.DatabaseKey,
	}
}

func (d *Dao) MakeOneSql4Insert(v interface{}, pk ...string) *mysql.OneSql {
	oneSql, _, _ := makeInsertSql(v, d.DatabaseKey, pk...)
	return oneSql
}

func (d *Dao) MakeOneSql4Update(v interface{}, pk ...string) *mysql.OneSql {
	oneSql, _, _ := makeUpdateSql(v, d.DatabaseKey, pk...)
	return oneSql

}

// ******************************  包装结束 ********************************
