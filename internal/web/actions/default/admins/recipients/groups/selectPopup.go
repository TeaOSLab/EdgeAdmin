package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) RunGet(params struct {
	GroupIds string
}) {
	selectedGroupIds := utils.SplitNumbers(params.GroupIds)

	groupsResp, err := this.RPC().MessageRecipientGroupRPC().FindAllEnabledMessageRecipientGroups(this.AdminContext(), &pb.FindAllEnabledMessageRecipientGroupsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	groupMaps := []maps.Map{}
	for _, group := range groupsResp.MessageRecipientGroups {
		if lists.ContainsInt64(selectedGroupIds, group.Id) {
			continue
		}
		groupMaps = append(groupMaps, maps.Map{
			"id":   group.Id,
			"name": group.Name,
			"isOn": group.IsOn,
		})
	}
	this.Data["groups"] = groupMaps

	this.Show()
}
