//
// Redis常用工具类
// 可通过配置文件启用集群模式还是主从模式
//
// 由于golang不支持泛型，对于获取数据需要调用显示指定
// 如果获取数据方法参数中的instance interface{}需要传入的是一个结构体类型，类似java的T.class。这里用：(*T)(nil)，string类型请使用nil
// 返回的都是interface{}类型的指针，遍历时直接可以强转：t := it.Next().(*T)
// 具体可以看redis_test.go测试Sample
//
// robot.guo

package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/exwallet/go-common/cache/redis/cluster"
	"github.com/exwallet/go-common/cache/redis/configuration"
	"github.com/exwallet/go-common/cache/redis/singleton"
	"github.com/exwallet/go-common/log"
	"github.com/exwallet/go-common/util/strutil"
	"reflect"
	"strconv"
	"sync"
)

var redisService RedisService
var lock sync.Mutex

type RedisService interface {
	InitRedis()
	Do(commandName string, args ...interface{}) (interface{}, error)
	Close() bool
}

func InitRedis(filepath string) {
	lock.Lock()
	defer lock.Unlock()
	if redisService != nil {
		return
	}
	redisConfig := configuration.InitDataSources(filepath)
	if redisConfig == nil {
		return
	}
	// 选择是集群模式还是主从模式
	if redisConfig.Mode == "cluster" {
		cluster := cluster.RedisCluster{}
		cluster.RedisConfig = redisConfig
		redisService = &cluster
	} else {
		singleton := singleton.RedisSingleton{}
		singleton.RedisConfig = redisConfig
		redisService = &singleton
	}
	// 初始化数据源
	redisService.InitRedis()
}

func Close() {
	if redisService != nil {
		redisService.Close()
	}
}

func Exists(k string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.Exists().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("EXISTS", k)
	if err != nil {
		log.Error("redis.Exists()出错了：%v\n", err)
		return false
	}
	if reply.(int64) > 0 {
		return true
	}
	return false
}

// Get*** 方法必须带回error , 如果redis挂了, 取到结果一直为空, 会影响限制类功能
func Get(k string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.Get().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("GET", k)
	if err != nil {
		log.Error("redis.Get()出错了：%v\n", err)
		return "", err
	}
	if reply == nil {
		return "", nil
	}
	return strutil.GetString(reply), nil
}

func Set(k string, v string) bool {
	return SetAndExpire(k, v, -1)
}

func SetAndExpire(k string, v string, seconds int) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.SetAndExpire().panic()：%v\n", err)
		}
	}()
	doFuc := func(k string, v string, seconds int) (interface{}, error) {
		if seconds >= 0 {
			return redisService.Do("SET", k, v, "EX", seconds)
		} else {
			return redisService.Do("SET", k, v)
		}
	}
	reply, err := doFuc(k, v, seconds)
	if err != nil {
		log.Error("redis.Set()出错了：%v\n", err)
		return false
	}
	if reply == "OK" {
		if seconds >= 0 {
			redisService.Do("EXPIRE", k, seconds)
		}
		return true
	}
	return false
}

func SetObj(k string, v interface{}) bool {
	return SetObjAndExpire(k, v, -1)
}

func SetObjAndExpire(k string, v interface{}, seconds int) bool {
	str, err := json.Marshal(v)
	if err != nil {
		return false
	}
	return SetAndExpire(k, string(str), seconds)
}

func SetList(k string, v *arraylist.List) bool {
	return SetListAndExpire(k, v, -1)
}

func SetListAndExpire(k string, v *arraylist.List, seconds int) bool {
	str, err := json.Marshal(v.Values())
	if err != nil {
		return false
	}
	return SetAndExpire(k, string(str), seconds)
}

func GetObj(k string, instance interface{}) (ret interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("redis.GetObj().panic()：%v", e))
		}
	}()
	str, e := Get(k)
	if e != nil {
		err = e
		return
	}
	if len(str) == 0 {
		return nil, nil
	}
	t := reflect.TypeOf(instance)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	ptr := reflect.New(t).Interface()
	//ptr := reflect.New(reflect.TypeOf(instance).Elem()).Interface()
	err = json.Unmarshal([]byte(str), ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func GetList(k string, instance interface{}) *arraylist.List {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.GetList().panic()：%v\n", err)
		}
	}()
	str, _ := Get(k)
	if len(str) == 0 {
		return nil
	}
	//
	instanceType := reflect.SliceOf(reflect.TypeOf(instance).Elem())
	elements := reflect.New(instanceType).Interface()
	err := json.Unmarshal([]byte(str), elements)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	//
	v := reflect.ValueOf(elements).Elem()
	// 反射出来就是reflect.Slice
	if v.Kind() != reflect.Slice {
		return nil
	}
	list := arraylist.New()
	for i := 0; i < v.Len(); i++ {
		list.Add(v.Index(i).Interface())
	}
	return list
}

func Delete(k string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.Delete().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("DEL", k)
	if err != nil {
		log.Error("redis.Delete()出错了：%v\n", err)
		return false
	}
	if reply.(int64) > 0 {
		return true
	}
	return false
}

func IncrBy(k string) int {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.IncrBy().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("INCRBY", k, 1)
	if err != nil {
		log.Error("redis.IncrBy()出错了：%v\n", err)
		return -1
	}
	v, err := strconv.Atoi(strutil.GetString(reply))
	return v
}

func LPush(k string, seconds int, datas ...interface{}) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.LPush().panic()：%v\n", err)
		}
	}()
	doFuc := func(k string, datas []interface{}) []interface{} {
		p := []interface{}{k}
		for _, v := range datas {
			vs := reflect.TypeOf(v)
			if vs.Kind() == reflect.Struct {
				bytes, err := json.Marshal(v)
				if err != nil {
					return nil
				}
				p = append(p, string(bytes))
			} else {
				p = append(p, v)
			}
		}
		return p
	}
	p := doFuc(k, datas)
	reply, err := redisService.Do("LPush", p...)
	if err != nil {
		log.Error("redis.LPush()出错了：%v\n", err)
		return false
	}
	if seconds >= 0 {
		redisService.Do("EXPIRE", k, seconds)
	}
	if reply.(int64) > 0 {
		return true
	}
	return false
}

func LRange(k string, start int, end int, instance interface{}) *arraylist.List {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.LRange().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("LRANGE", k, start, end)
	if err != nil {
		log.Error("redis.LRANGE()出错了：%v\n", err)
		return nil
	}
	arr := reply.([]interface{})
	l := arraylist.New()
	for _, v := range arr {
		str := strutil.GetString(v)
		if instance == nil {
			l.Add(str)
		} else {
			ptr := reflect.New(reflect.TypeOf(instance).Elem()).Interface()
			err := json.Unmarshal([]byte(str), ptr)
			if err != nil {
				return nil
			}
			l.Add(ptr)
		}
	}
	return l
}

func LSet(k string, index int, v interface{}, seconds int) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.LSet().panic()：%v\n", err)
		}
	}()
	p := v
	vs := reflect.TypeOf(v)
	if vs.Kind() == reflect.Struct {
		bytes, err := json.Marshal(v)
		if err != nil {
			return false
		}
		p = string(bytes)
	}
	reply, err := redisService.Do("LSET", k, index, p)
	if err != nil {
		log.Error("redis.LSET()出错了：%v\n", err)
		return false
	}
	if seconds >= 0 {
		redisService.Do("EXPIRE", k, seconds)
	}
	if reply == "OK" {
		return true
	}
	return false
}

func RPop(k string, instance interface{}) interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.RPop().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("RPOP", k)
	if err != nil {
		log.Error("redis.RPop()出错了：%v\n", err)
		return nil
	}
	str := strutil.GetString(reply)
	if instance == nil {
		return str
	} else {
		ptr := reflect.New(reflect.TypeOf(instance).Elem()).Interface()
		err := json.Unmarshal([]byte(str), ptr)
		if err != nil {
			return nil
		}
		return ptr
	}
}

func LLength(k string) int {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.LLength().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("LLEN", k)
	if err != nil {
		log.Error("redis.LLength()出错了：%v\n", err)
		return -1
	}
	v, err := strconv.Atoi(strutil.GetString(reply))
	return v
}

func HSet(k string, itemKey string, v interface{}, seconds int) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.HSet().panic()：%v\n", err)
		}
	}()
	p := v
	vs := reflect.TypeOf(v)
	if vs.Kind() == reflect.Struct {
		bytes, err := json.Marshal(v)
		if err != nil {
			return false
		}
		p = string(bytes)
	}
	_, err := redisService.Do("HSET", k, itemKey, p)
	if err != nil {
		log.Error("redis.HSet()出错了：%v\n", err)
		return false
	}
	if seconds >= 0 {
		redisService.Do("EXPIRE", k, seconds)
	}
	return true
}

func HGet(k string, itemKey string, instance interface{}) interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.HGet().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("HGET", k, itemKey)
	if err != nil {
		log.Error("redis.HGet()出错了：%v\n", err)
		return nil
	}
	str := strutil.GetString(reply)
	if instance == nil {
		return str
	} else {
		ptr := reflect.New(reflect.TypeOf(instance).Elem()).Interface()
		err := json.Unmarshal([]byte(str), ptr)
		if err != nil {
			return nil
		}
		return ptr
	}
}

func HDelete(k string, itemKey string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.HDelete().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("HDEL", k, itemKey)
	if err != nil {
		log.Error("redis.HDelete()出错了：%v\n", err)
		return false
	}
	if reply.(int64) > 0 {
		return true
	}
	return false
}

func HExists(k string, itemKey string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.HExists().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("HEXISTS", k, itemKey)
	if err != nil {
		log.Error("redis.HExists()出错了：%v\n", err)
		return false
	}
	if reply.(int64) > 0 {
		return true
	}
	return false
}

func HGetAll(k string, instance interface{}) map[string]interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Error("redis.HGetAll().panic()：%v\n", err)
		}
	}()
	reply, err := redisService.Do("HGETALL", k)
	if err != nil {
		log.Error("redis.HGetAll()出错了：%v\n", err)
		return nil
	}
	results := make(map[string]interface{})
	arr := reply.([]interface{})
	for i, v := range arr {
		if i%2 != 0 {
			continue
		}
		key := strutil.GetString(v)
		val := strutil.GetString(arr[i+1])
		if instance == nil {
			results[key] = val
		} else {
			ptr := reflect.New(reflect.TypeOf(instance).Elem()).Interface()
			err := json.Unmarshal([]byte(val), ptr)
			if err != nil {
				return nil
			}
			results[key] = ptr
		}
	}
	return results
}
