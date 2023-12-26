// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build linux

package executils

import (
	"golang.org/x/sys/unix"
	"io/fs"
	"os"
	"os/exec"
	"syscall"
)

// LookPath customize our LookPath() function, to work in broken $PATH environment variable
func LookPath(file string) (string, error) {
	result, err := exec.LookPath(file)
	if err == nil && len(result) > 0 {
		return result, nil
	}

	// add common dirs contains executable files these may be excluded in $PATH environment variable
	var binPaths = []string{
		"/usr/sbin",
		"/usr/bin",
		"/usr/local/sbin",
		"/usr/local/bin",
	}

	for _, binPath := range binPaths {
		var fullPath = binPath + string(os.PathSeparator) + file

		stat, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		if stat.IsDir() {
			return "", syscall.EISDIR
		}

		var mode = stat.Mode()
		if mode.IsDir() {
			return "", syscall.EISDIR
		}
		err = syscall.Faccessat(unix.AT_FDCWD, fullPath, unix.X_OK, unix.AT_EACCESS)
		if err == nil || (err != syscall.ENOSYS && err != syscall.EPERM) {
			return fullPath, err
		}
		if mode&0111 != 0 {
			return fullPath, nil
		}
		return "", fs.ErrPermission
	}

	return "", &exec.Error{
		Name: file,
		Err:  exec.ErrNotFound,
	}
}
