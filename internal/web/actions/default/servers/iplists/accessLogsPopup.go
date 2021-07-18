// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

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

	accessLogsResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
		Day:     timeutil.Format("Ymd"),
		Keyword: "ip:" + item.IpFrom + "," + item.IpTo,
		Size:    10,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var accessLogs = accessLogsResp.HttpAccessLogs
	if len(accessLogs) == 0 {
		accessLogs = []*pb.HTTPAccessLog{}
	}
	this.Data["accessLogs"] = accessLogs

	this.Show()
}
