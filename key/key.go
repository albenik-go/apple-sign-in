package key

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/pkg/errors"
)

var ErrPEMBlockMissing = errors.New("PEM block missing")

func ParsePrivateFromPEM(b []byte) (interface{}, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, ErrPEMBlockMissing
	}
	return x509.ParsePKCS8PrivateKey(block.Bytes)
}

func ReadPrivateFromPEMFile(file string) (interface{}, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "key file read error")
	}
	return ParsePrivateFromPEM(data)
}
