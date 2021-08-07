package recipients

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "recipient")
}

func (this *IndexAction) RunGet(params struct {
}) {
	// TODO 增加系统用户、媒介类型等条件搜索
	countResp, err := this.RPC().MessageRecipientRPC().CountAllEnabledMessageRecipients(this.AdminContext(), &pb.CountAllEnabledMessageRecipientsRequest{
		AdminId:                 0,
		MediaType:               "",
		MessageRecipientGroupId: 0,
		Keyword:                 "",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	recipientsResp, err := this.RPC().MessageRecipientRPC().ListEnabledMessageRecipients(this.AdminContext(), &pb.ListEnabledMessageRecipientsRequest{
		AdminId:                 0,
		MediaType:               "",
		MessageRecipientGroupId: 0,
		Keyword:                 "",
		Offset:                  page.Offset,
		Size:                    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	recipientMaps := []maps.Map{}
	for _, recipient := range recipientsResp.MessageRecipients {
		if recipient.Admin == nil {
			continue
		}
		if recipient.MessageMediaInstance == nil {
			continue
		}
		recipientMaps = append(recipientMaps, maps.Map{
			"id": recipient.Id,
			"admin": maps.Map{
				"id":       recipient.Admin.Id,
				"fullname": recipient.Admin.Fullname,
				"username": recipient.Admin.Username,
			},
			"groups": recipient.MessageRecipientGroups,
			"isOn":   recipient.IsOn,
			"instance": maps.Map{
				"name": recipient.MessageMediaInstance.Name,
			},
			"user":        recipient.User,
			"description": recipient.Description,
		})
	}
	this.Data["recipients"] = recipientMaps

	this.Show()
}
