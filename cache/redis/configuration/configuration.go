package configuration

import (
	"encoding/json"
	"github.com/exwallet/go-common/log"
	"github.com/exwallet/go-common/util/configuration/json"
	"time"
)

//配置信息
type JSONConfig struct {
	Mode        string
	Hosts       []string
	Ports       []int
	Pwd         string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

const (
	ConnTimeout  = time.Second
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 5 * time.Second
)

func InitDataSources(filePath string) *JSONConfig {
	str := configuration.ReadJSONFile(filePath)
	if len(str) == 0 {
		return nil
	}
	redisConfig := new(JSONConfig)
	err := json.Unmarshal([]byte(str), redisConfig)
	if err != nil {
		log.Error("无法解析的数据库配置文件：%s\n", str)
		log.Error("错误原因为：%v\n", err)
		return nil
	}
	//log.Debug("%+v\n", *redisConfig)
	return redisConfig
}
