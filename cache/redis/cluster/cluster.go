//
// redis集群模式接口实现
//
// robot.guo

package cluster

import (
	"errors"
	"fmt"
	"github.com/chasex/redis-go-cluster"
	"github.com/exwallet/go-common/cache/redis/configuration"
	"github.com/exwallet/go-common/gologger"
	"sync"
	"time"
)

type RedisCluster struct {
	init        bool
	RedisConfig *configuration.JSONConfig
	cluster     *redis.Cluster
	lock        sync.Mutex
}

func (r *RedisCluster) InitRedis() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.cluster != nil && r.init {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			gologger.Error("无法建立redis连接池，配置信息为:%+v\n", *r.RedisConfig)
			gologger.Error("错误信息为：%v\n", err)
			r.init = false
			return
		}
	}()
	nodes := make([]string, 0)
	for i, v := range r.RedisConfig.Hosts {
		nodes = append(nodes, fmt.Sprintf("%v:%v", v, r.RedisConfig.Ports[i]))
	}

	newCluster, err := redis.NewCluster(
		&redis.Options{
			StartNodes:   nodes,
			ConnTimeout:  configuration.ConnTimeout,
			ReadTimeout:  configuration.ReadTimeout,
			WriteTimeout: configuration.WriteTimeout,
			KeepAlive:    r.RedisConfig.MaxActive,
			AliveTime:    time.Duration(r.RedisConfig.IdleTimeout) * time.Second,
		})
	if err != nil {
		panic(err)
	}
	r.cluster = newCluster
	r.init = true
}

func (r *RedisCluster) Do(commandName string, args ...interface{}) (interface{}, error) {
	if !r.init || r.cluster == nil {
		return nil, errors.New("Redis集群未初始化")
	}
	return r.cluster.Do(commandName, args...)
}

func (r *RedisCluster) Close() bool {
	r.cluster.Close()
	return true
}
