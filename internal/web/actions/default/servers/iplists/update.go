// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/iplists/iplistutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	ListId int64
}) {
	err := InitIPList(this.Parent(), params.ListId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	ListId      int64
	Name        string
	Code        string
	Type        string
	Description string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.IPList_LogUpdateIPList, params.ListId)

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
		if listIdResp.IpListId > 0 && listIdResp.IpListId != params.ListId {
			this.FailField("code", "代号'"+params.Code+"'已经被别的名单占用，请更换一个")
			return
		}
	}

	_, err := this.RPC().IPListRPC().UpdateIPList(this.AdminContext(), &pb.UpdateIPListRequest{
		IpListId:    params.ListId,
		Name:        params.Name,
		Code:        params.Code,
		TimeoutJSON: nil,
		Description: params.Description,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
