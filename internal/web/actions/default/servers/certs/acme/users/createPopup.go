package users

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

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Email       string
	Description string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("email", params.Email).
		Require("请输入邮箱").
		Email("请输入正确的邮箱格式")

	createResp, err := this.RPC().ACMEUserRPC().CreateACMEUser(this.AdminContext(), &pb.CreateACMEUserRequest{
		Email:       params.Email,
		Description: params.Description,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 日志
	defer this.CreateLogInfo("创建ACME用户 %d", createResp.AcmeUserId)

	this.Success()
}
