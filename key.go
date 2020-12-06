package applesignin

import (
	"crypto/x509"
	"encoding/pem"
)

func ParsePrivateKey(b []byte) (interface{}, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, ErrNoPemBlockFound
	}
	return x509.ParsePKCS8PrivateKey(block.Bytes)
}
