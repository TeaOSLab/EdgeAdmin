package recipients

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type RecipientAction struct {
	actionutils.ParentAction
}

func (this *RecipientAction) Init() {
	this.Nav("", "", "recipient")
}

func (this *RecipientAction) RunGet(params struct {
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
			"id":          recipient.MessageMediaInstance.Id,
			"name":        recipient.MessageMediaInstance.Name,
			"description": recipient.MessageMediaInstance.Description,
		},
		"user":        recipient.User,
		"description": recipient.Description,
		"timeFrom":    recipient.TimeFrom,
		"timeTo":      recipient.TimeTo,
	}

	this.Show()
}
