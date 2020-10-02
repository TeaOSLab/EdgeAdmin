package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	t.Log(New("hello"))
	t.Log(Wrap(errors.New("hello")))
	t.Log(testError1())
	t.Log(Wrap(testError1()))
	t.Log(Wrap(testError2()))
}

func testError1() error {
	return New("test error1")
}

func testError2() error {
	return Wrap(testError1())
}
