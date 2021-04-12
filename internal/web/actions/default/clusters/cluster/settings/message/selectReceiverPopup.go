package message

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type SelectReceiverPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectReceiverPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectReceiverPopupAction) RunGet(params struct {
	RecipientIds string
	GroupIds     string
}) {
	recipientIds := utils.SplitNumbers(params.RecipientIds)
	groupIds := utils.SplitNumbers(params.GroupIds)

	// 所有接收人
	recipientsResp, err := this.RPC().MessageRecipientRPC().ListEnabledMessageRecipients(this.AdminContext(), &pb.ListEnabledMessageRecipientsRequest{
		AdminId:                 0,
		MediaType:               "",
		MessageRecipientGroupId: 0,
		Keyword:                 "",
		Offset:                  0,
		Size:                    1000, // TODO 支持搜索
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	recipientMaps := []maps.Map{}
	for _, recipient := range recipientsResp.MessageRecipients {
		if !recipient.IsOn {
			continue
		}
		if lists.ContainsInt64(recipientIds, recipient.Id) {
			continue
		}
		recipientMaps = append(recipientMaps, maps.Map{
			"id":           recipient.Id,
			"name":         recipient.Admin.Fullname,
			"instanceName": recipient.MessageMediaInstance.Name,
		})
	}
	this.Data["recipients"] = recipientMaps

	// 所有分组
	groupsResp, err := this.RPC().MessageRecipientGroupRPC().FindAllEnabledMessageRecipientGroups(this.AdminContext(), &pb.FindAllEnabledMessageRecipientGroupsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	groupMaps := []maps.Map{}
	for _, group := range groupsResp.MessageRecipientGroups {
		if !group.IsOn {
			continue
		}
		if lists.ContainsInt64(groupIds, group.Id) {
			continue
		}
		groupMaps = append(groupMaps, maps.Map{
			"id":   group.Id,
			"name": group.Name,
		})
	}
	this.Data["groups"] = groupMaps

	this.Show()
}
