// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"net"
)

// CheckLocalIPAction 检查IP是否为局域网IP
type CheckLocalIPAction struct {
	actionutils.ParentAction
}

func (this *CheckLocalIPAction) RunPost(params struct {
	Host string
}) {
	var ip = net.ParseIP(params.Host)
	if ip == nil {
		// 默认为true
		this.Data["isLocal"] = true
		this.Success()
	}

	var ipObj = ip.To4()
	if ipObj == nil {
		ipObj = ip.To16()
	}
	if ipObj == nil {
		// 默认为true
		this.Data["isLocal"] = true
		this.Success()
	}

	this.Data["isLocal"] = utils.IsLocalIP(ipObj)

	this.Success()
}
