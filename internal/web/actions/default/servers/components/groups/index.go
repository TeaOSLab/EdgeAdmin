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
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct{}) {
	groupsResp, err := this.RPC().ServerGroupRPC().FindAllEnabledServerGroups(this.AdminContext(), &pb.FindAllEnabledServerGroupsRequest{
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	groupMaps := []maps.Map{}
	for _, group := range groupsResp.ServerGroups {
		countResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithServerGroupId(this.AdminContext(), &pb.CountAllEnabledServersWithServerGroupIdRequest{ServerGroupId: group.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countServers := countResp.Count

		groupMaps = append(groupMaps, maps.Map{
			"id":           group.Id,
			"name":         group.Name,
			"countServers": countServers,
		})
	}
	this.Data["groups"] = groupMaps

	this.Show()
}
