// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package executils

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Cmd struct {
	name string
	args []string
	env  []string
	dir  string

	ctx        context.Context
	timeout    time.Duration
	cancelFunc func()

	captureStdout bool
	captureStderr bool

	stdout *bytes.Buffer
	stderr *bytes.Buffer

	rawCmd *exec.Cmd
}

func NewCmd(name string, args ...string) *Cmd {
	return &Cmd{
		name: name,
		args: args,
	}
}

func NewTimeoutCmd(timeout time.Duration, name string, args ...string) *Cmd {
	return (&Cmd{
		name: name,
		args: args,
	}).WithTimeout(timeout)
}

func (this *Cmd) WithTimeout(timeout time.Duration) *Cmd {
	this.timeout = timeout

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	this.ctx = ctx
	this.cancelFunc = cancelFunc

	return this
}

func (this *Cmd) WithStdout() *Cmd {
	this.captureStdout = true
	return this
}

func (this *Cmd) WithStderr() *Cmd {
	this.captureStderr = true
	return this
}

func (this *Cmd) WithEnv(env []string) *Cmd {
	this.env = env
	return this
}

func (this *Cmd) WithDir(dir string) *Cmd {
	this.dir = dir
	return this
}

func (this *Cmd) Start() error {
	var cmd = this.compose()
	return cmd.Start()
}

func (this *Cmd) Wait() error {
	var cmd = this.compose()
	return cmd.Wait()
}

func (this *Cmd) Run() error {
	if this.cancelFunc != nil {
		defer this.cancelFunc()
	}

	var cmd = this.compose()
	return cmd.Run()
}

func (this *Cmd) RawStdout() string {
	if this.stdout != nil {
		return this.stdout.String()
	}
	return ""
}

func (this *Cmd) Stdout() string {
	return strings.TrimSpace(this.RawStdout())
}

func (this *Cmd) RawStderr() string {
	if this.stderr != nil {
		return this.stderr.String()
	}
	return ""
}

func (this *Cmd) Stderr() string {
	return strings.TrimSpace(this.RawStderr())
}

func (this *Cmd) String() string {
	if this.rawCmd != nil {
		return this.rawCmd.String()
	}
	var newCmd = exec.Command(this.name, this.args...)
	return newCmd.String()
}

func (this *Cmd) Process() *os.Process {
	if this.rawCmd != nil {
		return this.rawCmd.Process
	}
	return nil
}

func (this *Cmd) compose() *exec.Cmd {
	if this.rawCmd != nil {
		return this.rawCmd
	}

	if this.ctx != nil {
		this.rawCmd = exec.CommandContext(this.ctx, this.name, this.args...)
	} else {
		this.rawCmd = exec.Command(this.name, this.args...)
	}

	if this.env != nil {
		this.rawCmd.Env = this.env
	}

	if len(this.dir) > 0 {
		this.rawCmd.Dir = this.dir
	}

	if this.captureStdout {
		this.stdout = &bytes.Buffer{}
		this.rawCmd.Stdout = this.stdout
	}
	if this.captureStderr {
		this.stderr = &bytes.Buffer{}
		this.rawCmd.Stderr = this.stderr
	}

	return this.rawCmd
}
