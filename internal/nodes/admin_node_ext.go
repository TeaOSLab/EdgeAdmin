// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package nodes

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
	"github.com/iwind/TeaGo/logs"
)

// 启动IP库
func (this *AdminNode) startIPLibrary() {
	logs.Println("[NODE]initializing ip library ...")
	err := iplibrary.InitDefault()
	if err != nil {
		logs.Println("[NODE]initialize ip library failed: "+err.Error())
	}
}
