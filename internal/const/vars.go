// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package teaconst

import (
	"os"
	"strings"
)

var (
	IsRecoverMode = false

	IsDemoMode         = false
	ErrorDemoOperation = "DEMO模式下无法进行创建、修改、删除等操作"

	NewVersionCode        = "" // 有新的版本
	NewVersionDownloadURL = "" // 新版本下载地址

	IsMain = checkMain()
)

// 检查是否为主程序
func checkMain() bool {
	if len(os.Args) == 1 ||
		(len(os.Args) >= 2 && os.Args[1] == "pprof") {
		return true
	}
	exe, _ := os.Executable()
	return strings.HasSuffix(exe, ".test") ||
		strings.HasSuffix(exe, ".test.exe") ||
		strings.Contains(exe, "___")
}
