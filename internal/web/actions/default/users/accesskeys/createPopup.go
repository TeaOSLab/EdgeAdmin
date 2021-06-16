// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package accesskeys

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	UserId int64
}) {
	this.Data["userId"] = params.UserId
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	UserId      int64
	Description string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("description", params.Description).
		Require("请输入备注")

	accessKeyIdResp, err := this.RPC().UserAccessKeyRPC().CreateUserAccessKey(this.AdminContext(), &pb.CreateUserAccessKeyRequest{
		UserId:      params.UserId,
		Description: params.Description,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	defer this.CreateLogInfo("创建AccessKey %d", accessKeyIdResp.UserAccessKeyId)

	this.Success()
}
