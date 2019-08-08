package gorsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
)

// MarshalPKCS8PrivateKey 私钥解析
func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {

	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}

	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, _ := asn1.Marshal(info)
	return k

}
