package applesignin

import (
	"fmt"

	"github.com/pkg/errors"
)

var ErrNoPemBlockFound = errors.New("no PEM block found")

// ErrorResponse see https://developer.apple.com/documentation/sign_in_with_apple/errorresponse.
type ErrorResponse struct {
	Reason string `json:"error"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("apple error response: %s", r.Reason)
}
