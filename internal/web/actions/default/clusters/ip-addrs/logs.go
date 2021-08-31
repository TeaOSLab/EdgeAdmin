// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipaddrs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type LogsAction struct {
	actionutils.ParentAction
}

func (this *LogsAction) Init() {
	this.Nav("", "", "log")
}

func (this *LogsAction) RunGet(params struct {
}) {
	countResp, err := this.RPC().NodeIPAddressLogRPC().CountAllNodeIPAddressLogs(this.AdminContext(), &pb.CountAllNodeIPAddressLogsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var page = this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	logsResp, err := this.RPC().NodeIPAddressLogRPC().ListNodeIPAddressLogs(this.AdminContext(), &pb.ListNodeIPAddressLogsRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var logMaps = []maps.Map{}
	for _, log := range logsResp.NodeIPAddressLogs {
		var adminMap maps.Map
		if log.Admin != nil {
			adminMap = maps.Map{
				"id":   log.Admin.Id,
				"name": log.Admin.Fullname,
			}
		} else {
			adminMap = maps.Map{
				"id":   0,
				"name": "[系统]",
			}
		}

		var addrMap maps.Map
		if log.NodeIPAddress != nil {
			var addr = log.NodeIPAddress
			if addr != nil {
				addrMap = maps.Map{
					"id": addr.Id,
					"ip": addr.Ip,
				}
			}
		}

		logMaps = append(logMaps, maps.Map{
			"id":          log.Id,
			"description": log.Description,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt),
			"isUp":        log.IsUp,
			"isOn":        log.IsOn,
			"canAccess":   log.CanAccess,
			"admin":       adminMap,
			"addr":        addrMap,
		})
	}
	this.Data["logs"] = logMaps

	this.Show()
}
