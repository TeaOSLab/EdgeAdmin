// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipbox

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct {
	Ip string
}) {
	this.Data["ip"] = params.Ip

	// IP信息
	regionResp, err := this.RPC().IPLibraryRPC().LookupIPRegion(this.AdminContext(), &pb.LookupIPRegionRequest{Ip: params.Ip})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if regionResp.IpRegion != nil {
		this.Data["regions"] = regionResp.IpRegion.Summary
	} else {
		this.Data["regions"] = ""
	}
	this.Data["isp"] = regionResp.IpRegion.Isp

	// IP列表
	ipListResp, err := this.RPC().IPListRPC().FindEnabledIPListContainsIP(this.AdminContext(), &pb.FindEnabledIPListContainsIPRequest{Ip: params.Ip})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var ipListMaps = []maps.Map{}
	for _, ipList := range ipListResp.IpLists {
		ipListMaps = append(ipListMaps, maps.Map{
			"id":   ipList.Id,
			"name": ipList.Name,
			"type": ipList.Type,
		})
	}
	this.Data["ipLists"] = ipListMaps

	// 所有公用的IP列表
	publicBlackIPListResp, err := this.RPC().IPListRPC().ListEnabledIPLists(this.AdminContext(), &pb.ListEnabledIPListsRequest{
		Type:     "black",
		IsPublic: true,
		Keyword:  "",
		Offset:   0,
		Size:     10, // TODO 将来考虑到支持更多的黑名单
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var publicBlackIPListMaps = []maps.Map{}
	for _, ipList := range publicBlackIPListResp.IpLists {
		publicBlackIPListMaps = append(publicBlackIPListMaps, maps.Map{
			"id":   ipList.Id,
			"name": ipList.Name,
			"type": ipList.Type,
		})
	}
	this.Data["publicBlackIPLists"] = publicBlackIPListMaps

	// 访问日志
	accessLogsResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
		Day:  timeutil.Format("Ymd"),
		Ip:   params.Ip,
		Size: 20,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var accessLogs = accessLogsResp.HttpAccessLogs
	if len(accessLogs) == 0 {
		// 查询昨天
		accessLogsResp, err = this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
			Day:  timeutil.Format("Ymd", time.Now().AddDate(0, 0, -1)),
			Ip:   params.Ip,
			Size: 20,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		accessLogs = accessLogsResp.HttpAccessLogs

		if len(accessLogs) == 0 {
			accessLogs = []*pb.HTTPAccessLog{}
		}
	}
	this.Data["accessLogs"] = accessLogs

	this.Show()
}
