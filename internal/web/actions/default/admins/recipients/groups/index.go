package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "group")
}

func (this *IndexAction) RunGet(params struct{}) {
	groupsResp, err := this.RPC().MessageRecipientGroupRPC().FindAllEnabledMessageRecipientGroups(this.AdminContext(), &pb.FindAllEnabledMessageRecipientGroupsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	groupMaps := []maps.Map{}
	for _, group := range groupsResp.MessageRecipientGroups {
		groupMaps = append(groupMaps, maps.Map{
			"id":   group.Id,
			"name": group.Name,
			"isOn": group.IsOn,
		})
	}
	this.Data["groups"] = groupMaps

	this.Show()
}
