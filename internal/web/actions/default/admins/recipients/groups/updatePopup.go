package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	GroupId int64
}) {
	groupResp, err := this.RPC().MessageRecipientGroupRPC().FindEnabledMessageRecipientGroup(this.AdminContext(), &pb.FindEnabledMessageRecipientGroupRequest{MessageRecipientGroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	group := groupResp.MessageRecipientGroup
	if group == nil {
		this.NotFound("messageRecipientGroup", params.GroupId)
		return
	}

	this.Data["group"] = maps.Map{
		"id":   group.Id,
		"name": group.Name,
		"isOn": group.IsOn,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	GroupId int64
	Name    string
	IsOn    bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")

	_, err := this.RPC().MessageRecipientGroupRPC().UpdateMessageRecipientGroup(this.AdminContext(), &pb.UpdateMessageRecipientGroupRequest{
		MessageRecipientGroupId: params.GroupId,
		Name:                    params.Name,
		IsOn:                    params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
