// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

import "io"

type ProgressWriter struct {
	rawWriter io.Writer
	total     int64
	written   int64
}

func NewProgressWriter(rawWriter io.Writer, total int64) *ProgressWriter {
	return &ProgressWriter{
		rawWriter: rawWriter,
		total:     total,
	}
}

func (this *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = this.rawWriter.Write(p)
	this.written += int64(n)
	return
}

func (this *ProgressWriter) Progress() float32 {
	if this.total <= 0 {
		return 0
	}
	return float32(float64(this.written) / float64(this.total))
}
