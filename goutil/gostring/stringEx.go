/*
 * @Author: kidd
 * @Date: 1/29/19 6:28 PM
 *
 */

package gostring

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/exwallet/go-common/goutil/goaes"
	"reflect"
	"strconv"
	"strings"
)

// String AT的字符串库
type String string

const (
	// NilString 定义空串的关键字
	NilString String = "<nil>"
)

//StringProtocal 强化String的扩展方法
type StringProtocal interface {
	String() string
	Length() int
	IsNil() bool
	IsEmpty() bool
}

// String 转string
func (s String) String() string {
	return string(s)
}

// Length 字符串长度，支持中文字等等的计算
func (s String) Length() int {
	length := len([]rune(s.String()))
	return length
}

// IsNil 判断是否空指针
func (s String) IsNil() bool {

	if s == NilString {
		return true
	}

	return false
}

// IsEmpty 是否空值
func (s String) IsEmpty() bool {

	if len(s) == 0 {
		return true
	}

	return false
}

// Int String类型转为int类型
func (s String) Int(def ...int) (int, bool) {
	var (
		value int
		err   error
	)
	if value, err = strconv.Atoi(s.String()); err != nil {
		if len(def) > 0 {
			return def[0], true
			//value = def[0]
		} else {
			return 0, false
		}
	}
	return value, true
}

// UInt8 String类型转为uint8类型
func (s String) UInt8(def ...int) (uint8, bool) {
	value, ok := s.UInt64(def...)
	return uint8(value), ok
}

// UInt16 String类型转为uint32类型
func (s String) UInt16(def ...int) (uint16, bool) {
	value, ok := s.UInt64(def...)
	return uint16(value), ok
}

// UInt32 String类型转为uint32类型
func (s String) UInt32(def ...int) (uint32, bool) {
	value, ok := s.UInt64(def...)
	return uint32(value), ok
}

// UInt64 String类型转为uint64类型
func (s String) UInt64(def ...int) (uint64, bool) {
	var (
		value uint64
		err   error
	)
	if value, err = strconv.ParseUint(s.String(), 10, 64); err != nil {
		if len(def) > 0 {
			//value = uint64(def[0])
			return uint64(def[0]), true
		}
		return 0, false
	}
	return value, true
}

// Int8 String类型转为int8类型
func (s String) Int8(def ...int) (int8, bool) {
	value, ok := s.Int64(def...)
	return int8(value), ok
}

// Int16 String类型转为int32类型
func (s String) Int16(def ...int) (int16, bool) {
	value, ok := s.Int64(def...)
	return int16(value), ok
}

// Int32 String类型转为int32类型
func (s String) Int32(def ...int) (int32, bool) {
	value, ok := s.Int64(def...)
	return int32(value), ok
}

// Int64 String类型转为int64类型
func (s String) Int64(def ...int) (int64, bool) {
	var (
		value int64
		err   error
	)
	if value, err = strconv.ParseInt(s.String(), 10, 64); err != nil {
		if len(def) > 0 {
			//value = int64(def[0])
			return int64(def[0]), true
		}
		return 0, false
	}
	return value, true
}

// Bool string转布尔型
func (s String) Bool(def ...bool) bool {
	if strings.ToLower(s.String()) == "true" {
		return true
	}
	if strings.ToLower(s.String()) == "false" {
		return false
	}
	i, ok := s.Int()
	if ok && i == 1 {
		return true
	}
	if ok && i == 0 {
		return false
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

// Float32 String转为float32
func (s String) Float32(def ...float32) (float32, bool) {
	value, ok := s.Float64()
	return float32(value), ok
}

// Float64 String转为float64
func (s String) Float64(def ...float64) (float64, bool) {
	var (
		value float64
		err   error
	)
	if value, err = strconv.ParseFloat(s.String(), 32); err != nil {
		if len(def) > 0 {
			value = def[0]
			return def[0], true
		}
		return 0, false
	}
	return value, true
}

// AES AES加密
// key 密钥key hex字符串
// return 密文
func (s String) AES(key string) (string, error) {
	var (
		plantext   []byte
		keybyte    []byte
		err        error
		ciphertext []byte
		result     string
	)
	plantext = []byte(s.String())
	if keybyte, err = hex.DecodeString(key); err != nil {
		return "", err
	}

	if ciphertext, err = goaes.AESEncrypt(plantext, keybyte); err != nil {
		return "", err
	}

	//转为base64
	result = base64.StdEncoding.EncodeToString(ciphertext)

	return result, err
}

// UnAES 通过base64编码的aes密文初始化一个字符串
// aesBase64string base64编码的aes密文
// key 密钥字符串
// return 明文
func (s *String) UnAES(aesBase64string string, key string) error {
	var (
		plantext   []byte
		keybyte    []byte
		err        error
		ciphertext []byte
		result     String
	)
	if keybyte, err = hex.DecodeString(key); err != nil {
		return err
	}
	if ciphertext, err = base64.StdEncoding.DecodeString(aesBase64string); err != nil {
		return err
	}
	if plantext, err = goaes.AESDecrypt(ciphertext, keybyte); err != nil {
		return err
	}
	result = String(plantext)
	*s = result

	return nil
}

// NewStringByInt 通过int初始化字符串
func NewStringByInt(v int64) String {
	str := strconv.FormatInt(v, 10)
	return String(str)
}

// NewStringByUInt 通过int初始化字符串
func NewStringByUInt(v uint64) String {
	str := strconv.FormatUint(v, 10)
	return String(str)
}

// NewStringByBool 通过bool初始化字符串
func NewStringByBool(v bool) String {
	str := strconv.FormatBool(v)
	return String(str)
}

// NewStringByFloat 通过float初始化字符串
func NewStringByFloat(v float64) String {
	str := strconv.FormatFloat(v, 'f', -1, 64)
	return String(str)
}

// NewString 初始化字符串，自动转型
func NewString(value interface{}, def ...String) String {

	val := reflect.ValueOf(value) //读取变量的值，可能是指针或值

	switch val.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewStringByInt(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewStringByUInt(val.Uint())
	case reflect.Float32, reflect.Float64:
		return NewStringByFloat(val.Float())
	case reflect.Bool:
		return NewStringByBool(val.Bool())
	case reflect.String:
		return String(val.String())
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct: //如果类型为数组，要继续迭代元素
		jsonstr, _ := json.Marshal(value)
		return String(jsonstr)
	default:
		if len(def) > 0 {
			return String(def[0])
		}

	}
	return ""
}

// MD5 字符串转为MD5后的hash hex
func (s String) MD5() string {
	return GetMD5(s.String())
	//return mdStr
}

// SHA1 字符串转为SHA1后的hash hex
func (s String) SHA1() string {
	return GetSHA1(s.String())
}

// SHA256 字符串转为SHA256后的hash hex
func (s String) SHA256() string {
	return GetSHA256(s.String())
}

// HmacSHA1 字符串转为HmacSHA1后的hash hex
func (s String) HmacSHA1(secret string) string {
	//hash := crypto.HmacSHA1(secret, []byte(s))
	//mdStr := hex.EncodeToString(hash)
	return GetHmacSHA1(secret, s.String())
}

// HmacMD5 字符串转为HmacMD5后的hash md5
func (s String) HmacMD5(secret string) string {
	//hash := crypto.HmacMD5(secret, []byte(s))
	//mdStr := hex.EncodeToString(hash)
	return GetHmacMD5(secret, s.String())
}

func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		panic("start is wrong")
	}
	if end < 0 || end > length {
		panic("end is wrong")
	}
	return string(rs[start:end:end])
}

func FormatStruct(v interface{}) string {
	objstr, _ := json.MarshalIndent(v, "", " ")
	return string(objstr) //log.Debugf("event:%v", string(objstr))
}
