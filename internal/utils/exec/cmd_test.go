// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package executils_test

import (
	executils "github.com/TeaOSLab/EdgeAdmin/internal/utils/exec"
	"testing"
	"time"
)

func TestNewTimeoutCmd_Sleep(t *testing.T) {
	var cmd = executils.NewTimeoutCmd(1*time.Second, "sleep", "3")
	cmd.WithStdout()
	cmd.WithStderr()
	err := cmd.Run()
	t.Log("error:", err)
	t.Log("stdout:", cmd.Stdout())
	t.Log("stderr:", cmd.Stderr())
}

func TestNewTimeoutCmd_Echo(t *testing.T) {
	var cmd = executils.NewTimeoutCmd(10*time.Second, "echo", "-n", "hello")
	cmd.WithStdout()
	cmd.WithStderr()
	err := cmd.Run()
	t.Log("error:", err)
	t.Log("stdout:", cmd.Stdout())
	t.Log("stderr:", cmd.Stderr())
}

func TestNewTimeoutCmd_Echo2(t *testing.T) {
	var cmd = executils.NewCmd("echo", "hello")
	cmd.WithStdout()
	cmd.WithStderr()
	err := cmd.Run()
	t.Log("error:", err)
	t.Log("stdout:", cmd.Stdout())
	t.Log("raw stdout:", cmd.RawStdout())
	t.Log("stderr:", cmd.Stderr())
	t.Log("raw stderr:", cmd.RawStderr())
}

func TestNewTimeoutCmd_Echo3(t *testing.T) {
	var cmd = executils.NewCmd("echo", "-n", "hello")
	err := cmd.Run()
	t.Log("error:", err)
	t.Log("stdout:", cmd.Stdout())
	t.Log("stderr:", cmd.Stderr())
}

func TestCmd_Process(t *testing.T) {
	var cmd = executils.NewCmd("echo", "-n", "hello")
	err := cmd.Run()
	t.Log("error:", err)
	t.Log(cmd.Process())
}

func TestNewTimeoutCmd_String(t *testing.T) {
	var cmd = executils.NewCmd("echo", "-n", "hello")
	t.Log("stdout:", cmd.String())
}
