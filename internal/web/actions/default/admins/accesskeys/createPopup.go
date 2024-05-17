// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package accesskeys

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
	AdminId int64
}) {
	this.Data["adminId"] = params.AdminId
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	AdminId     int64
	Description string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("description", params.Description).
		Require("请输入备注")

	accessKeyIdResp, err := this.RPC().UserAccessKeyRPC().CreateUserAccessKey(this.AdminContext(), &pb.CreateUserAccessKeyRequest{
		AdminId:     params.AdminId,
		Description: params.Description,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	defer this.CreateLogInfo(codes.UserAccessKey_LogCreateUserAccessKey, accessKeyIdResp.UserAccessKeyId)

	this.Success()
}
