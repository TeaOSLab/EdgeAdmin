// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package apinodeutils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// DeployFile 部署文件描述
type DeployFile struct {
	OS      string
	Arch    string
	Version string
	Path    string
}

// Sum 计算概要
func (this *DeployFile) Sum() (string, error) {
	fp, err := os.Open(this.Path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = fp.Close()
	}()

	m := md5.New()
	buffer := make([]byte, 128*1024)
	for {
		n, err := fp.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		_, err = m.Write(buffer[:n])
		if err != nil {
			return "", err
		}
	}
	sum := m.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}

// Read 读取一个片段数据
func (this *DeployFile) Read(offset int64) (data []byte, newOffset int64, err error) {
	fp, err := os.Open(this.Path)
	if err != nil {
		return nil, offset, err
	}
	defer func() {
		_ = fp.Close()
	}()

	stat, err := fp.Stat()
	if err != nil {
		return nil, offset, err
	}
	if offset >= stat.Size() {
		return nil, offset, io.EOF
	}

	_, err = fp.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, offset, err
	}

	buffer := make([]byte, 128*1024)
	n, err := fp.Read(buffer)
	if err != nil {
		return nil, offset, err
	}

	return buffer[:n], offset + int64(n), nil
}
