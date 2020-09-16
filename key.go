package applesignin

import (
	"crypto/x509"
	"encoding/pem"
)

func LoadP8CertByByte(str []byte) (interface{}, error) {
	block, _ := pem.Decode(str)
	return x509.ParsePKCS8PrivateKey(block.Bytes)
}
