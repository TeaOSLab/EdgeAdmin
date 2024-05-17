// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/iplists/iplistutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
	Code        string
	Type        string
	Description string
	IsGlobal    bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var listId int64 = 0
	defer func() {
		defer this.CreateLogInfo(codes.IPList_LogCreateIPList, listId)
	}()

	params.Must.
		Field("name", params.Name).
		Require("请输入名称")

	if len(params.Code) > 0 {
		if !iplistutils.ValidateIPListCode(params.Code) {
			this.FailField("code", "代号格式错误，只能是英文字母、数字、中划线、下划线的组合")
			return
		}

		listIdResp, err := this.RPC().IPListRPC().FindIPListIdWithCode(this.AdminContext(), &pb.FindIPListIdWithCodeRequest{Code: params.Code})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if listIdResp.IpListId > 0 {
			this.FailField("code", "代号'"+params.Code+"'已经被别的名单占用，请更换一个")
			return
		}
	}

	createResp, err := this.RPC().IPListRPC().CreateIPList(this.AdminContext(), &pb.CreateIPListRequest{
		Type:        params.Type,
		Name:        params.Name,
		Code:        params.Code,
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
