// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

var SharedLogger = NewLogger()

type Logger struct {
	c chan string
}

func NewLogger() *Logger {
	return &Logger{
		c: make(chan string, 1024),
	}
}

func (this *Logger) Push(msg string) {
	select {
	case this.c <- msg:
	default:

	}
}

func (this *Logger) ReadAll() (msgList []string) {
	msgList = []string{}

	for {
		select {
		case msg := <-this.c:
			msgList = append(msgList, msg)
		default:
			return
		}
	}
}

func (this *Logger) Reset() {
	for {
		select {
		case <-this.c:
		default:
			return
		}
	}
}
