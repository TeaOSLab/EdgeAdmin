// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type BindHTTPFirewallPopupAction struct {
	actionutils.ParentAction
}

func (this *BindHTTPFirewallPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *BindHTTPFirewallPopupAction) RunGet(params struct {
	HttpFirewallPolicyId int64
	Type                 string
}) {
	this.Data["httpFirewallPolicyId"] = params.HttpFirewallPolicyId

	// 获取已经选中的名单IDs
	var selectedIds = []int64{}
	inboundConfig, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyInboundConfig(this.AdminContext(), params.HttpFirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if inboundConfig != nil {
		for _, ref := range inboundConfig.PublicAllowListRefs {
			selectedIds = append(selectedIds, ref.ListId)
		}
		for _, ref := range inboundConfig.PublicDenyListRefs {
			selectedIds = append(selectedIds, ref.ListId)
		}
		for _, ref := range inboundConfig.PublicGreyListRefs {
			selectedIds = append(selectedIds, ref.ListId)
		}
	}

	// 公共的名单
	countResp, err := this.RPC().IPListRPC().CountAllEnabledIPLists(this.AdminContext(), &pb.CountAllEnabledIPListsRequest{
		Type:     params.Type,
		IsPublic: true,
		Keyword:  "",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	listsResp, err := this.RPC().IPListRPC().ListEnabledIPLists(this.AdminContext(), &pb.ListEnabledIPListsRequest{
		Type:     params.Type,
		IsPublic: true,
		Keyword:  "",
		Offset:   page.Offset,
		Size:     page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var listMaps = []maps.Map{}
	for _, list := range listsResp.IpLists {
		// 包含的IP数量
		countItemsResp, err := this.RPC().IPItemRPC().CountIPItemsWithListId(this.AdminContext(), &pb.CountIPItemsWithListIdRequest{IpListId: list.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var countItems = countItemsResp.Count

		listMaps = append(listMaps, maps.Map{
			"id":          list.Id,
			"isOn":        list.IsOn,
			"name":        list.Name,
			"description": list.Description,
			"countItems":  countItems,
			"type":        list.Type,
			"isSelected":  lists.ContainsInt64(selectedIds, list.Id),
		})
	}
	this.Data["lists"] = listMaps

	this.Show()
}

func (this *BindHTTPFirewallPopupAction) RunPost(params struct {
	HttpFirewallPolicyId int64
	ListId               int64

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.IPList_LogBindIPListWAFPolicy, params.ListId, params.HttpFirewallPolicyId)

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
	inboundConfig.AddPublicList(list.Id, list.Type)

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
