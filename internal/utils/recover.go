package utils

import (
	"runtime/debug"
)

func Recover() {
	e := recover()
	if e != nil {
		debug.PrintStack()
	}
}
