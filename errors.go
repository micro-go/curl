package curl

import (
	"errors"
)

var (
	badRequestErr     = errors.New("Bad request")
	noResponseBodyErr = errors.New("No response body")
)

func mergeError(a error, b error) error {
	if a != nil {
		return a
	}
	return b
}
