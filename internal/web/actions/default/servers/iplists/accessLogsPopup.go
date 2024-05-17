// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type AccessLogsPopupAction struct {
	actionutils.ParentAction
}

func (this *AccessLogsPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *AccessLogsPopupAction) RunGet(params struct {
	ItemId int64
}) {
	itemResp, err := this.RPC().IPItemRPC().FindEnabledIPItem(this.AdminContext(), &pb.FindEnabledIPItemRequest{IpItemId: params.ItemId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var item = itemResp.IpItem
	if item == nil {
		this.NotFound("ipItem", params.ItemId)
		return
	}

	this.Data["ipFrom"] = item.IpFrom
	this.Data["ipTo"] = item.IpTo

	// 多找几个Partition
	var day = timeutil.FormatTime("Ymd", item.CreatedAt)
	partitionsResp, err := this.RPC().HTTPAccessLogRPC().FindHTTPAccessLogPartitions(this.AdminContext(), &pb.FindHTTPAccessLogPartitionsRequest{Day: day})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var hasAccessLogs = false

	for _, partition := range partitionsResp.ReversePartitions {
		accessLogsResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
			Partition: partition,
			Day:       day,
			Keyword:   "ip:" + item.IpFrom + "," + item.IpTo,
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
			break
		}
	}

	if !hasAccessLogs {
		this.Data["accessLogs"] = []interface{}{}
	}

	this.Show()
}
