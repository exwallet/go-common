/*
 * @Author: kidd
 * @Date: 1/19/19 10:24 PM
 */

package googleAuth

import (
	"fmt"
	"testing"
	"time"
)

func TestGA(t *testing.T) {
	// 生成种子
	seed := _GenerateSeed()
	fmt.Println("种子: ", seed)
	// 通过种子生成密钥
	key, _ := _GenerateSecretKey(seed)
	fmt.Println("密钥: ", key)

	// 通过密钥+时间生成验证码
	rs := _GetNewCode(key, time.Now().Unix())
	fmt.Println("验证码: ", rs)
	fmt.Println("开始睡眠延迟中,请耐心等待...")
	time.Sleep(5 * time.Second)
	// 校验已有验证码
	fmt.Println("校验结果: ", ValidCode(key, rs))
}
