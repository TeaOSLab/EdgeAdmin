// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"net"
	"strings"
)

type ValidateIPsAction struct {
	actionutils.ParentAction
}

func (this *ValidateIPsAction) RunPost(params struct {
	Ips string
}) {
	var ips = params.Ips
	if len(ips) == 0 {
		this.Data["ips"] = []string{}
		this.Success()
	}

	var ipSlice = strings.Split(ips, "\n")
	var result = []string{}
	for _, ip := range ipSlice {
		ip = strings.TrimSpace(ip)
		if len(ip) == 0 {
			continue
		}
		data := net.ParseIP(ip)
		if len(data) == 0 {
			this.Data["failIP"] = ip
			this.Fail()
		}
		result = append(result, ip)
	}

	this.Data["ips"] = result
	this.Success()
}
