package groups

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
	Name string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")

	_, err := this.RPC().MessageRecipientGroupRPC().CreateMessageRecipientGroup(this.AdminContext(), &pb.CreateMessageRecipientGroupRequest{Name: params.Name})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
