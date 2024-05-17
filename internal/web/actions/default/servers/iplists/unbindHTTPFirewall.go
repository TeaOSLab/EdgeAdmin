// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
)

type UnbindHTTPFirewallAction struct {
	actionutils.ParentAction
}

func (this *UnbindHTTPFirewallAction) RunPost(params struct {
	HttpFirewallPolicyId int64
	ListId               int64
}) {
	defer this.CreateLogInfo(codes.IPList_LogUnbindIPListWAFPolicy, params.ListId, params.HttpFirewallPolicyId)

	// List类型
	listResp, err := this.RPC().IPListRPC().FindEnabledIPList(this.AdminContext(), &pb.FindEnabledIPListRequest{IpListId: params.ListId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var list = listResp.IpList
	if list == nil {
		this.Fail("找不到要使用的IP名单")
	}

	// 已经绑定的
	inboundConfig, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyInboundConfig(this.AdminContext(), params.HttpFirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if inboundConfig == nil {
		inboundConfig = &firewallconfigs.HTTPFirewallInboundConfig{IsOn: true}
	}
	inboundConfig.RemovePublicList(list.Id, list.Type)

	inboundJSON, err := json.Marshal(inboundConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().HTTPFirewallPolicyRPC().UpdateHTTPFirewallInboundConfig(this.AdminContext(), &pb.UpdateHTTPFirewallInboundConfigRequest{
		HttpFirewallPolicyId: params.HttpFirewallPolicyId,
		InboundJSON:          inboundJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
