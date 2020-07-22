// +build windows

package apps

import (
	"errors"
	"os"
)

// lock file
func LockFile(fp *os.File) error {
	return errors.New("not implemented on windows")
}

func UnlockFile(fp *os.File) error {
	return errors.New("not implemented on windows")
}
