// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package ipbox

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
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


	this.Data["regions"] = ""
	this.Data["isp"] = ""

	var ipRegion = iplibrary.LookupIP(params.Ip)
	if ipRegion != nil && ipRegion.IsOk() {
		this.Data["regions"] = ipRegion.RegionSummary()
		this.Data["isp"] = ipRegion.ProviderName()
	}

	// IP列表
	ipListResp, err := this.RPC().IPListRPC().FindEnabledIPListContainsIP(this.AdminContext(), &pb.FindEnabledIPListContainsIPRequest{
		Ip: params.Ip,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var ipListMaps = []maps.Map{}
	for _, ipList := range ipListResp.IpLists {
		itemsResp, err := this.RPC().IPItemRPC().ListIPItemsWithListId(this.AdminContext(), &pb.ListIPItemsWithListIdRequest{
			IpListId: ipList.Id,
			Keyword:  "",
			IpFrom:   params.Ip,
			IpTo:     "",
			Offset:   0,
			Size:     1,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var items = itemsResp.IpItems
		if len(items) == 0 {
			continue
		}
		var item = items[0]

		var expiredTime = ""
		if item.ExpiredAt > 0 {
			expiredTime = timeutil.FormatTime("Y-m-d H:i:s", item.ExpiredAt)
		}

		ipListMaps = append(ipListMaps, maps.Map{
			"id":              ipList.Id,
			"name":            ipList.Name,
			"type":            ipList.Type,
			"itemExpiredTime": expiredTime,
			"itemId":          item.Id,
		})
	}
	this.Data["ipLists"] = ipListMaps

	// 所有公用的IP列表
	publicBlackIPListResp, err := this.RPC().IPListRPC().ListEnabledIPLists(this.AdminContext(), &pb.ListEnabledIPListsRequest{
		Type:     "black",
		IsPublic: true,
		Keyword:  "",
		Offset:   0,
		Size:     20, // TODO 将来考虑到支持更多的黑名单
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
	var hasAccessLogs = false
Loop:
	for _, day := range []string{timeutil.Format("Ymd"), timeutil.Format("Ymd", time.Now().AddDate(0, 0, -1))} {
		partitionsResp, err := this.RPC().HTTPAccessLogRPC().FindHTTPAccessLogPartitions(this.AdminContext(), &pb.FindHTTPAccessLogPartitionsRequest{Day: day})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, partition := range partitionsResp.ReversePartitions {
			accessLogsResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
				Partition: partition,
				Day:       day,
				Ip:        params.Ip,
				Size:      20,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			var accessLogs = accessLogsResp.HttpAccessLogs
			if len(accessLogs) > 0 {
				this.Data["accessLogs"] = accessLogs
				hasAccessLogs = true
				break Loop
			}
		}
	}

	if !hasAccessLogs {
		this.Data["accessLogs"] = []interface{}{}
	}

	this.Show()
}
