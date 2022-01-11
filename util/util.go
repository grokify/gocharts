package util

import (
	errors "errors"
	"fmt"

	errors2 "github.com/pkg/errors"
)

func main() {
	err1 := errors.New("foo")
	err2 := errors2.Wrap(err1, "bar")
	fmt.Println(err2.Error())
	fmt.Println(ErrorWrap(err1, "bar"))
}

func ErrorWrap(err error, prefix string) error {
	return errors.New(prefix + ": " + err.Error())
}
