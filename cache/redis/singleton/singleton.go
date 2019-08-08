//
// 单个实例，也就是主从模式接口实现
//
// robot.guo
package singleton

import (
	"fmt"
	"github.com/exwallet/go-common/cache/redis/configuration"
	"github.com/exwallet/go-common/gologger"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"sync"
	"time"
)

type RedisSingleton struct {
	init        bool
	RedisConfig *configuration.JSONConfig
	redisPool   *redis.Pool
	lock        sync.Mutex
}

func (r *RedisSingleton) InitRedis() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.redisPool != nil && r.init {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			gologger.Error("无法建立redis连接池，配置信息为:%+v\n", *r.RedisConfig)
			gologger.Error("错误信息为：%v\n", err)
		}
	}()
	// redis pool
	r.redisPool = &redis.Pool{
		MaxIdle:     r.RedisConfig.MaxIdle,
		MaxActive:   r.RedisConfig.MaxActive,
		IdleTimeout: time.Duration(r.RedisConfig.IdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", r.RedisConfig.Hosts[0], strconv.Itoa(r.RedisConfig.Ports[0])),
				redis.DialPassword(r.RedisConfig.Pwd),
				//redis.DialDatabase(0),
				redis.DialConnectTimeout(configuration.ConnTimeout),
				redis.DialReadTimeout(configuration.ReadTimeout),
				redis.DialWriteTimeout(configuration.WriteTimeout))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
	r.init = true
}

func (r *RedisSingleton) borrowResource() redis.Conn {
	if !r.init || r.redisPool == nil {
		return nil
	}
	return r.redisPool.Get()
}

func (r *RedisSingleton) returnResource(rc redis.Conn) {
	if err := rc.Close(); err != nil {
		gologger.Error("释放Redis连接出错了\n")
	}
}

func (r *RedisSingleton) Do(commandName string, args ...interface{}) (interface{}, error) {
	rc := r.borrowResource()
	//Err returns a non-nil value when the connection is not usable.
	if err := rc.Err(); err != nil {
		gologger.Error("redis borrowResource is not usable：%v\n", err)
		return nil, err
	}
	defer r.returnResource(rc)
	return rc.Do(commandName, args...)
}

func (r *RedisSingleton) Close() bool {
	if r.redisPool != nil {
		err := r.redisPool.Close()
		if err != nil {
			return false
		}
	}
	return true
}
