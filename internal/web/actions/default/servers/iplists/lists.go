// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/ipconfigs"
	"github.com/iwind/TeaGo/maps"
)

type ListsAction struct {
	actionutils.ParentAction
}

func (this *ListsAction) Init() {
	this.Nav("", "", "lists")
}

func (this *ListsAction) RunGet(params struct {
	Type    string
	Keyword string
}) {
	if len(params.Type) == 0 {
		params.Type = ipconfigs.IPListTypeBlack
	}
	this.Data["type"] = params.Type
	this.Data["keyword"] = params.Keyword

	countResp, err := this.RPC().IPListRPC().CountAllEnabledIPLists(this.AdminContext(), &pb.CountAllEnabledIPListsRequest{
		Type:     params.Type,
		IsPublic: true,
		Keyword:  params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var page = this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	listsResp, err := this.RPC().IPListRPC().ListEnabledIPLists(this.AdminContext(), &pb.ListEnabledIPListsRequest{
		Type:     params.Type,
		IsPublic: true,
		Keyword:  params.Keyword,
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
			"code":        list.Code,
			"description": list.Description,
			"countItems":  countItems,
			"type":        list.Type,
			"isGlobal":    list.IsGlobal,
		})
	}
	this.Data["lists"] = listMaps

	this.Show()
}
