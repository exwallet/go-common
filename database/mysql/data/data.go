//
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
	"bytes"
	"database/sql"
	"fmt"
	"github.com/exwallet/go-common/database/mysql"
	"github.com/exwallet/go-common/gologger"
	"github.com/exwallet/go-common/util/strutil"
	"reflect"
	"strconv"
	"strings"
)

type CallWithMysqlTx interface {
	service() bool
}

const (
	TagPrimaryKey = "pk"
)

// 取主键字段名
func getPrimaryKey(ts reflect.Type, vs reflect.Value, pk ...string) (pkName string, pkVal interface{}) {
	var _pk string
	if len(pk) > 0 {
		_pk = pk[0]
	} else {
		for i := 0; i < ts.NumField(); i++ {
			fi := ts.Field(i)
			if fi.Tag.Get(TagPrimaryKey) != "" {
				_pk = fi.Name
				break
			}
		}
	}
	if _pk == "" {
		return
	}
	val := vs.FieldByName(_pk)
	if val.Kind() == reflect.Struct {
		return
	}
	switch val.Kind() {
	case reflect.Int64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int:
		pkName = _pk
		pkVal = val.Int()
		return
	case reflect.String:
		pkName = _pk
		pkVal = val.String()
		return
	}
	return
}

// 执行事务操作
func execute(txs map[string]*sql.Tx, sqls []*mysql.OneSql) (success bool) {
	if sqls == nil || len(sqls) == 0 {
		return
	}
	var myTx *sql.Tx
	// 必须捕获panic，并回滚事务
	defer func() {
		if err := recover(); err != nil {
			rollback(txs)
			success = false
			gologger.Error("捕获到panic抛回异常并回滚事务，异常信息：%v", err)
			return
		}
	}()
	//it := sqls.Iterator()
	for _, sql := range sqls {
		//获取连接
		if tx, ok := txs[sql.DatabaseKey]; ok {
			myTx = tx
		} else {
			db := mysql.GetConnection(sql.DatabaseKey)
			// 没有连接源
			if db == nil {
				gologger.Error("databaseKey(%s)，找不到数据源，导致事务回滚", sql.DatabaseKey)
				rollback(txs)
				return
			}
			newTx, err := db.Begin()
			if err != nil {
				gologger.Error("databaseKey(%s)，事务开启失败，导致事务回滚", sql.DatabaseKey)
				rollback(txs)
				return
			}
			txs[sql.DatabaseKey] = newTx
			myTx = newTx
		}
		// Prepare
		stmt, err := myTx.Prepare(sql.Sql)
		if err != nil {
			gologger.Error("execute tx.Prepare(): databaseKey(%s)，%s，%v，导致事务回滚", sql.DatabaseKey, sql.Sql, sql.Params)
			gologger.Error("error: %s", err.Error())
			rollback(txs)
			return
		}
		// execute
		res, err := stmt.Exec(sql.Params...)
		if err != nil {
			gologger.Error("execute stmt.Exec(): databaseKey(%s)，%s，%v，导致事务回滚\n", sql.DatabaseKey, sql.Sql, sql.Params)
			gologger.Error("error: %s", err.Error())
			rollback(txs)
			return
		}
		// RowsAffected
		num, err := res.RowsAffected()
		if err != nil {
			gologger.Error("execute res.RowsAffected(): databaseKey(%s)，%s，%v，导致事务回滚", sql.DatabaseKey, sql.Sql, sql.Params)
			gologger.Error("error: %s", err.Error())
			rollback(txs)
			return
		}
		// judgement
		if !sql.FutureEfRows(int(num)) {
			gologger.Error("execute sql.FutureEfRows(): databaseKey(%s)，%s，%v，与预期行数不一，预期影响：%d行，实际影响：%d行，导致事务回滚\n", sql.DatabaseKey, sql.Sql, sql.Params, sql.EfRows, num)
			rollback(txs)
			return
		}
		// close statement
		stmt.Close()
	}
	// execute success
	success = true
	return
}

func rollback(txs map[string]*sql.Tx) {
	for k, v := range txs {
		err := v.Rollback()
		if err != nil {
			gologger.Error("execute tx.rollback(): databaseKey(%s)，事务回滚失败了，有可能导致数据不同步\n", k)
			// rollback(txs)
			return
		}
	}
}

func makeUpdateSql(v interface{}, databaseKey string, pk ...string) (sql *mysql.OneSql, table string, err error) {
	vs := reflect.Indirect(reflect.ValueOf(v))
	ts := vs.Type()
	//
	table = ts.Name()
	pkName, pkVal := getPrimaryKey(ts, vs, pk...)
	if pkName == "" {
		//err = fmt.Errorf("%s表找不到主键", table)
		pkName = "id"
		return
	}
	fields := make([]string, 0, ts.NumField())
	for i := 0; i < ts.NumField(); i++ {
		fi := ts.Field(i)
		if fi.Name == pkName {
			continue
		}
		if !strutil.IsPrevUpper(fi.Name) {
			continue
		}
		fields = append(fields, fi.Name)
	}
	//
	var (
		results bytes.Buffer
		_fs     []string
		params  []interface{}
	)
	results.WriteString("update " + ts.Name() + " set ")
	for _, f := range fields {
		val := vs.FieldByName(f)
		if val.Kind() != reflect.Struct {
			_fs = append(_fs, "`"+f+"`=?")
			switch val.Kind() {
			case reflect.Int64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int:
				params = append(params, val.Int())
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				params = append(params, val.Uint())
			case reflect.Float64, reflect.Float32:
				params = append(params, val.Float())
			case reflect.String:
				params = append(params, val.String())
			default:
				params = append(params, val.Interface())
			}
		}
	}
	if len(params) == 0 {
		err = fmt.Errorf("%s表没有有效字段", ts.Name())
		return
	}

	results.WriteString(strings.Join(_fs, ","))
	results.WriteString(" where " + pkName + "=?")
	params = append(params, pkVal)
	//
	gologger.Debug("--> sql: %s \n--> params: %q\n", results.String(), params)
	gologger.Debug("--> params个数 :%d\n", len(params))
	sql = &mysql.OneSql{results.String(), 1, params, databaseKey}
	return
}

func makeInsertSql(v interface{}, databaseKey string, pk ...string) (sql *mysql.OneSql, table string, err error) {
	vs := reflect.Indirect(reflect.ValueOf(v))
	ts := vs.Type()
	//
	table = ts.Name()
	pkName, _ := getPrimaryKey(ts, vs, pk...)
	//if pkName == "" {
	//	err = fmt.Errorf("%s表找不到主键", ts.Name())
	//	return
	//}
	fields := make([]string, 0, ts.NumField())
	for i := 0; i < ts.NumField(); i++ {
		fi := ts.Field(i)
		// 去掉自增字段
		if fi.Name == pkName {
			continue
		}
		fields = append(fields, fi.Name)
	}

	var (
		results bytes.Buffer
		_fs     []string
		_val    []string
		params  []interface{}
	)
	//var results, b1, b2 bytes.Buffer
	results.WriteString("insert into " + ts.Name() + "(")

	for _, f := range fields {
		val := vs.FieldByName(f)
		if val.Kind() != reflect.Struct {
			_fs = append(_fs, "`"+f+"`")
			_val = append(_val, "?")
			switch val.Kind() {
			case reflect.Int64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int:
				params = append(params, val.Int())
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				params = append(params, val.Uint())
			case reflect.Float64, reflect.Float32:
				params = append(params, val.Float())
			case reflect.String:
				params = append(params, val.String())
			default:
				params = append(params, val.Interface())
			}
		}
	}
	if len(params) == 0 {
		err = fmt.Errorf("%s表没有有效字段", ts.Name())
		return
	}
	//
	results.WriteString(strings.Join(_fs, ","))
	results.WriteString(") values(")
	results.WriteString(strings.Join(_val, ","))
	results.WriteString(")")
	//
	//logger.Debug("--> sql: %s \n--> params: %q\n", results.String(), params)
	//logger.Debug("--> params个数 :%d\n", len(params))
	sql = &mysql.OneSql{results.String(), 1, params, databaseKey}
	return
}

func mapping(columns []string, values []sql.RawBytes, instance interface{}) interface{} {
	//
	rowMap := make(map[string]string)
	for i, col := range values {
		if col != nil {
			rowMap[strings.ToLower(columns[i])] = string(col)
		}
	}
	// 自动判断是否指针类
	v := reflect.New(reflect.Indirect(reflect.ValueOf(instance)).Type()).Interface()
	//// 确定是struct类型
	//v := reflect.New(reflect.TypeOf(instance)).Interface()
	//// 确定是指针类
	//v := reflect.New(reflect.TypeOf(instance).Elem()).Interface()
	// instance 指针
	//v := reflect.New(reflect.TypeOf(instance).Elem()).Interface()
	// reflect.Type

	// instance 指针
	//v := reflect.New(reflect.TypeOf(instance).Elem()).Interface()
	// reflect.Type
	ts := reflect.TypeOf(v).Elem()
	// reflect.Value
	vs := reflect.ValueOf(v).Elem()
	for i := 0; i < ts.NumField(); i++ {
		fi := ts.Field(i)

		if strutil.IsPrevUpper(fi.Name) && fi.Type.Kind() != reflect.Struct {
			if mapVal, ok := rowMap[strings.ToLower(fi.Name)]; ok {
				switch fi.Type.Kind() {
				case reflect.Int64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int:
					if intVal, err := strconv.ParseInt(mapVal, 10, 64); err == nil {
						vs.FieldByName(fi.Name).SetInt(intVal)
					}
				case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
					if uIntVal, err := strconv.ParseUint(mapVal, 10, 64); err == nil {
						vs.FieldByName(fi.Name).SetUint(uIntVal)
					}
				case reflect.Float64, reflect.Float32:
					if floatValue, err := strconv.ParseFloat(mapVal, 64); err == nil {
						vs.FieldByName(fi.Name).SetFloat(floatValue)
					}
				case reflect.String:
					vs.FieldByName(fi.Name).SetString(mapVal)
				default:

				}
			}
		}
	}
	return v
}
