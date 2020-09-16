package applesignin

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func readResponseBody(r *http.Response) ([]byte, error) {
	if r.ContentLength <= 0 || r.ContentLength > 1024 {
		return nil, errors.Errorf("unexpected response body length %d", r.ContentLength)
	}

	// TODO think about sync.Pool
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Errorf("cannot read response body")
	}
	return buf, nil
}

func readResponseBodyText(r *http.Response) (s string, err error) {
	var b []byte
	if b, err = readResponseBody(r); err == nil {
		s = string(b)
	}
	return
}
