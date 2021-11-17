// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	Type string
}) {
	this.Data["type"] = params.Type

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name        string
	Type        string
	Description string
	IsGlobal    bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var listId int64 = 0
	defer func() {
		defer this.CreateLogInfo("创建IP名单 %d", listId)
	}()

	params.Must.
		Field("name", params.Name).
		Require("请输入名称")

	createResp, err := this.RPC().IPListRPC().CreateIPList(this.AdminContext(), &pb.CreateIPListRequest{
		Type:        params.Type,
		Name:        params.Name,
		Code:        "",
		TimeoutJSON: nil,
		IsPublic:    true,
		Description: params.Description,
		IsGlobal:    params.IsGlobal,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	listId = createResp.IpListId

	this.Data["list"] = maps.Map{
		"type": params.Type,
	}

	this.Success()
}
