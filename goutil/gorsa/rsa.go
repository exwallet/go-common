/*
 * @Author: kidd
 * @Date: 1/19/19 10:24 PM
 */

package gorsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
	"runtime"
)

const (

	// RSAAlgorithmSign RSA签名算法
	RSAAlgorithmSign = crypto.SHA256
)

type RsaKey struct {
	PrivateKey []byte
	PublicKey  []byte
}

// string -> hex string
func (k *RsaKey) Encrypt(plainText string) (string, error) {
	if b, e := Encrypt([]byte(plainText), k.PublicKey); e == nil {
		return base64.StdEncoding.EncodeToString(b), nil
		//return fmt.Sprintf("%x", b), nil
	} else {
		return "", e
	}
}

// hex string -> bytes
func (k *RsaKey) Decrypt(plainText string) (string, error) {
	decodeString, e := base64.StdEncoding.DecodeString(plainText)
	//bytes, e := hex.DecodeString(plainText)
	if e != nil {
		return "", e
	}
	b, err := Decrypt(decodeString, k.PrivateKey)
	if err == nil {
		return string(b), nil
	}
	return "", err
}

func (k *RsaKey) Sign(msg string) (string, error) {
	cryptText, err := Sign(msg, k.PrivateKey)
	if err != nil {
		return "", err
	}
	return cryptText, nil
}

func (k *RsaKey) VerifySign(data string, sign string) bool {
	return VerifySign(data, sign, k.PublicKey)
}

/*
	非对称加密需要生成一对密钥而不是一个密钥，所以这里加密之前需要获取一对密钥分别为公钥和私钥
	一次性生成公钥和私钥
		加密:  明文的E	次方 Mod N  输出密文
		解密:  密文的D    次方 Mod N  输出明文
		加密操作需要消耗很长的时间 ? 加密速度会快
		数据加密之后不能被轻易的解密出来
*/

func NewRsaKey(bits ...int) (key *RsaKey, err error) {
	key = &RsaKey{}
	//1. GetKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥。
	//Reader是一个全局、共享的密码用强随机数生成器。在Unix类型系统中，会从/dev/urandom读取；
	// 而Windows中会调用CryptGenRandom API。
	var privateKey *rsa.PrivateKey
	var e error
	if len(bits) > 0 {
		privateKey, e = rsa.GenerateKey(rand.Reader, bits[0])
	} else {
		privateKey, e = rsa.GenerateKey(rand.Reader, 2048)
	}
	if e != nil {
		err = e
		return
	}
	// 将公钥和私钥持久的保存下来, 将这些内容保存到文件中

	//2.x509标准 按照一定标准的标准对数据进行格式化.序列化.编码
	// MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码。
	x509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//derStream := MarshalPKCS8PrivateKey(privateKey)

	//3.使用pem格式对x509输出的内容进行编码 base64编码--> 64个字符 0-9 a-z A-Z + /  总共64个字符
	// 构建一个block结构体
	privateBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509PrivateKey,
	}
	key.PrivateKey = pem.EncodeToMemory(privateBlock)

	//=================================生成公钥===============================================
	//x509对公钥进行编码,有很多不可见的字符,乱码 MarshalPKIXPublicKey将公钥序列化为PKIX格式DER编码。
	// 注意,这里传入的必须是指针
	x509PublicKey, e := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if e != nil {
		err = e
		return
	}
	//3.2 构建一个block对象
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509PublicKey,
	}
	//3.3 使用pem编码
	//pem.Encode(publicFile,&publicBlock)
	key.PublicKey = pem.EncodeToMemory(block)
	return
}

func Encrypt(plainText, key []byte) (cryptText []byte, err error) {
	//1. pem 解码
	block, _ := pem.Decode(key)
	//防止用户传的密钥不正确导致panic,这里恢复程序并打印错误
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Println("runtime err:", err, "请检查密钥是否正确")
			default:
				log.Println("error:", err)
			}
		}
	}()
	//2. block中的Bytes是x509编码的内容, x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return []byte{}, err //出错返回错误
	}
	//3.1 类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	//4. 使用公钥对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return []byte{}, err //出错返回错误
	}
	return cipherText, nil
}

func Decrypt(cryptText, key []byte) (plainText []byte, err error) {
	//1. pem格式解码
	block, _ := pem.Decode(key)
	//防止用户传的密钥不正确导致panic,这里恢复程序并打印错误
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Println("runtime err:", err, "请检查密钥是否正确")
			default:
				log.Println("error:", err)
			}
		}
	}()
	//2.x509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return []byte{}, err
	}
	//3. 解密操作
	plainText, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, cryptText)
	if err != nil {
		return []byte{}, err
	}
	return plainText, nil
}

//实现的是利用RSA数字签名的函数，注意：用公钥加密，私钥解密就是加密通信，用私钥加密，公钥验证相当于数字签名
func Sign(msg string, Key []byte) (cryptText string, err error) {
	//1. pem格式解码
	block, _ := pem.Decode(Key)
	//防止用户传的密钥不正确导致panic,这里恢复程序并打印错误
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Println("runtime err:", err, "请检查密钥是否正确")
			default:
				log.Println("error:", err)
			}
		}
	}()
	//2.x509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	// 计算消息的hash值
	myHash := RSAAlgorithmSign.New()
	myHash.Write([]byte(msg))
	hashed := myHash.Sum(nil)

	//SignPKCS1v15使用RSA PKCS#1 v1.5规定的RSASSA-PKCS1-V1_5-SIGN签名方案计算签名。注意hashed必须是使用提供给本函数的hash参数对（要签名的）原始数据进行hash的结果。
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, RSAAlgorithmSign, hashed)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sign), nil //返回签名后的消息

}

//验证签名，验证签名用公钥验证，如果可以解密验证说明签名正确，否则错误
func VerifySign(data string, sign string, Key []byte) bool { //如果解密正确，那么就返回true,否着返回false
	//1. pem格式解码
	block, _ := pem.Decode(Key)
	//防止用户传的密钥不正确导致panic,这里恢复程序并打印错误
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Println("runtime err:", err, "请检查密钥是否正确")
			default:
				log.Println("error:", err)
			}
		}
	}()
	//2.x509解码
	publicInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey := publicInterface.(*rsa.PublicKey)

	h := RSAAlgorithmSign.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
	if err != nil {
		return false
	}
	return rsa.VerifyPKCS1v15(publicKey, RSAAlgorithmSign, hashed, decodedSign) == nil

}
