//
// mysql配置信息及全局数据源
//
// robot.guo

package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/exwallet/go-common/log"
	"github.com/exwallet/go-common/util/configuration/json"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

var mysqlConfig *MysqlConfig
var lock sync.Mutex

type MysqlConfig struct {
	dataSources map[string]*sql.DB
	jsonConfigs map[string]JSONConfig
	lock        sync.Mutex
}

//配置信息
type JSONConfig struct {
	Key             string
	Ip              string
	User            string
	Pwd             string
	Db              string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

func InitDataSources(filepath string) {
	lock.Lock()
	defer lock.Unlock()

	if mysqlConfig != nil {
		return
	}

	mysqlConfig = new(MysqlConfig)
	mysqlConfig.jsonConfigs = make(map[string]JSONConfig)
	mysqlConfig.dataSources = make(map[string]*sql.DB)

	str := configuration.ReadJSONFile(filepath)
	if len(str) == 0 {
		return
	}
	fmt.Println(str)

	arr := make([]JSONConfig, 0, 8)
	err := json.Unmarshal([]byte(str), &arr)
	if err != nil {
		log.Debug("无法解析的数据库配置文件：%s\n", str)
		log.Debug("错误原因为：%v\n", err)
		// 错误的配置文件
		panic(err)
	}

	for _, jsonConfig := range arr {
		//已经初始化过了
		if _, ok := mysqlConfig.jsonConfigs[jsonConfig.Key]; ok {
			return
		}
		mysqlConfig.jsonConfigs[jsonConfig.Key] = jsonConfig
	}
}

func GetConnection(databaseKey string) *sql.DB {
	mysqlConfig.lock.Lock()
	defer mysqlConfig.lock.Unlock()
	// 如果已经建立了连接
	if v, ok := mysqlConfig.dataSources[databaseKey]; ok {
		return v
	}
	// 通过配置建立mysql连接源
	jsonConfig, ok := mysqlConfig.jsonConfigs[databaseKey]
	if !ok {
		return nil
	}
	return openConnection(&jsonConfig)
}

func openConnection(jsonConfig *JSONConfig) *sql.DB {
	// 建立连接
	url := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&allowOldPasswords=1", jsonConfig.User, jsonConfig.Pwd, jsonConfig.Ip, jsonConfig.Db)
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Error("连接mysql失败，databaseKey: %v，原因:%v\n", jsonConfig.Key, err)
		return nil
	}
	// 设置连接池最大链接数--不能大于数据库设置的最大链接数
	db.SetMaxOpenConns(jsonConfig.MaxOpenConns)
	// 设置最大空闲链接数--小于设置的链接数
	db.SetMaxIdleConns(jsonConfig.MaxIdleConns)
	// 设置数据库链接超时时间--不能大于数据库设置的超时时间
	db.SetConnMaxLifetime(time.Duration(jsonConfig.ConnMaxLifetime) * time.Minute)
	// 保存当前数据库连接
	mysqlConfig.dataSources[jsonConfig.Key] = db
	log.Info("mysql配置信息user[%s] db[%s] ip[%s]:%+v\n", jsonConfig.User, jsonConfig.Db, jsonConfig.Ip)
	return db
}

func Close() {
	if mysqlConfig == nil {
		return
	}
	for k, db := range mysqlConfig.dataSources {
		// 连接已经关闭了
		if err := db.Ping(); err != nil {
			continue
		}
		if err := db.Close(); err == nil {
			delete(mysqlConfig.dataSources, k)
			log.Error("关闭mysql连接成功：%v\n", k)
		} else {
			log.Error("关闭mysql连接失败：%v\n", err)
		}
	}
}
