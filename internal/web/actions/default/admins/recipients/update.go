package recipients

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	RecipientId int64
}) {
	recipientResp, err := this.RPC().MessageRecipientRPC().FindEnabledMessageRecipient(this.AdminContext(), &pb.FindEnabledMessageRecipientRequest{MessageRecipientId: params.RecipientId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	recipient := recipientResp.MessageRecipient
	if recipient == nil || recipient.Admin == nil || recipient.MessageMediaInstance == nil {
		this.NotFound("messageRecipient", params.RecipientId)
		return
	}

	this.Data["recipient"] = maps.Map{
		"id": recipient.Id,
		"admin": maps.Map{
			"id":       recipient.Admin.Id,
			"fullname": recipient.Admin.Fullname,
			"username": recipient.Admin.Username,
		},
		"groups": recipient.MessageRecipientGroups,
		"isOn":   recipient.IsOn,
		"instance": maps.Map{
			"id":   recipient.MessageMediaInstance.Id,
			"name": recipient.MessageMediaInstance.Name,
		},
		"user":        recipient.User,
		"description": recipient.Description,
	}

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	RecipientId int64

	AdminId    int64
	User       string
	InstanceId int64

	GroupIds    string
	Description string
	IsOn        bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改媒介接收人 %d", params.RecipientId)

	params.Must.
		Field("adminId", params.AdminId).
		Gt(0, "请选择系统用户").
		Field("instanceId", params.InstanceId).
		Gt(0, "请选择媒介")

	groupIds := utils.SplitNumbers(params.GroupIds)

	_, err := this.RPC().MessageRecipientRPC().UpdateMessageRecipient(this.AdminContext(), &pb.UpdateMessageRecipientRequest{
		MessageRecipientId: params.RecipientId,
		AdminId:            params.AdminId,
		InstanceId:         params.InstanceId,
		User:               params.User,
		GroupIds:           groupIds,
		Description:        params.Description,
		IsOn:               params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
