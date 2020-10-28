package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type OptionsAction struct {
	actionutils.ParentAction
}

func (this *OptionsAction) RunPost(params struct {
	ClusterId int64
}) {
	groupsResp, err := this.RPC().NodeGroupRPC().FindAllEnabledNodeGroupsWithClusterId(this.AdminContext(), &pb.FindAllEnabledNodeGroupsWithClusterIdRequest{ClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
	}

	groupMaps := []maps.Map{}
	for _, group := range groupsResp.Groups {
		groupMaps = append(groupMaps, maps.Map{
			"id":   group.Id,
			"name": group.Name,
		})
	}
	this.Data["groups"] = groupMaps

	this.Success()
}
