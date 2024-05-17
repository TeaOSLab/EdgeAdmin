// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type TestAction struct {
	actionutils.ParentAction
}

func (this *TestAction) Init() {
	this.Nav("", "", "test")
}

func (this *TestAction) RunGet(params struct {
	ListId int64
}) {
	err := InitIPList(this.Parent(), params.ListId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}

func (this *TestAction) RunPost(params struct {
	ListId int64
	Ip     string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	resp, err := this.RPC().IPItemRPC().CheckIPItemStatus(this.AdminContext(), &pb.CheckIPItemStatusRequest{
		IpListId: params.ListId,
		Ip:       params.Ip,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var resultMap = maps.Map{
		"isDone":    true,
		"isFound":   resp.IsFound,
		"isOk":      resp.IsOk,
		"error":     resp.Error,
		"isAllowed": resp.IsAllowed,
	}

	if resp.IpItem != nil {
		resultMap["item"] = maps.Map{
			"id":             resp.IpItem.Id,
			"value":          resp.IpItem.Value,
			"ipFrom":         resp.IpItem.IpFrom,
			"ipTo":           resp.IpItem.IpTo,
			"reason":         resp.IpItem.Reason,
			"expiredAt":      resp.IpItem.ExpiredAt,
			"createdTime":    timeutil.FormatTime("Y-m-d", resp.IpItem.CreatedAt),
			"expiredTime":    timeutil.FormatTime("Y-m-d H:i:s", resp.IpItem.ExpiredAt),
			"type":           resp.IpItem.Type,
			"eventLevelName": firewallconfigs.FindFirewallEventLevelName(resp.IpItem.EventLevel),
			"listType":       resp.IpItem.ListType,
		}
	}

	this.Data["result"] = resultMap

	this.Success()
}
