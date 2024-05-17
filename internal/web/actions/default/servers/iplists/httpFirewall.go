// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/ipconfigs"
	"github.com/iwind/TeaGo/maps"
)

// HttpFirewallAction 显示已经绑定的IP名单
type HttpFirewallAction struct {
	actionutils.ParentAction
}

func (this *HttpFirewallAction) RunPost(params struct {
	HttpFirewallPolicyId int64
	Type                 string
}) {
	inboundConfig, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyInboundConfig(this.AdminContext(), params.HttpFirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if inboundConfig == nil {
		inboundConfig = &firewallconfigs.HTTPFirewallInboundConfig{IsOn: true}
	}
	var refs []*ipconfigs.IPListRef
	switch params.Type {
	case ipconfigs.IPListTypeBlack:
		refs = inboundConfig.PublicDenyListRefs
	case ipconfigs.IPListTypeWhite:
		refs = inboundConfig.PublicAllowListRefs
	case ipconfigs.IPListTypeGrey:
		refs = inboundConfig.PublicGreyListRefs
	}

	var listMaps = []maps.Map{}
	for _, ref := range refs {
		listResp, err := this.RPC().IPListRPC().FindEnabledIPList(this.AdminContext(), &pb.FindEnabledIPListRequest{IpListId: ref.ListId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var list = listResp.IpList
		if list == nil {
			continue
		}

		listMaps = append(listMaps, maps.Map{
			"id":   list.Id,
			"name": list.Name,
		})
	}
	this.Data["lists"] = listMaps

	this.Success()
}
