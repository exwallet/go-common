/*
 * @Author: kidd
 * @Date: 1/24/19 1:34 PM
 *
 */

package bucketFilter

import (
	"fmt"
	"github.com/exwallet/go-common/cache/redis"
	"testing"
	"time"
)

func Test(t *testing.T) {
	redis.InitRedis("conf/redis.json")
	// 建立个限制桶, 单位时间10秒, 单位时间内最大尝试次数2次,  允许重试间隔4秒
	b := NewBucketFilter("test", "a13a_", 10, 4, 2)
	key1 := "192.168.1.1"

	for {
		<-time.After(time.Second * 1)
		e := b.CheckAvailable(key1)
		if e == nil {
			b.Success(key1)
			fmt.Println("成功")
		}

	}

}
