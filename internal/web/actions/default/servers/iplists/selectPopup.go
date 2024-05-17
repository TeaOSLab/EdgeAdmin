// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectPopupAction) RunGet(params struct {
	Type string
}) {
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
		})
	}
	this.Data["lists"] = listMaps

	this.Show()
}
