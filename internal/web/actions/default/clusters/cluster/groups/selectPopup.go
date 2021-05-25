package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectPopupAction) RunGet(params struct {
	ClusterId int64
}) {
	groupsResp, err := this.RPC().NodeGroupRPC().FindAllEnabledNodeGroupsWithNodeClusterId(this.AdminContext(), &pb.FindAllEnabledNodeGroupsWithNodeClusterIdRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
	}

	groupMaps := []maps.Map{}
	for _, group := range groupsResp.NodeGroups {
		groupMaps = append(groupMaps, maps.Map{
			"id":   group.Id,
			"name": group.Name,
		})
	}
	this.Data["groups"] = groupMaps

	this.Show()
}

func (this *SelectPopupAction) RunPost(params struct {
	GroupId int64

	Must *actions.Must
}) {
	if params.GroupId <= 0 {
		this.Fail("请选择要使用的分组")
	}

	groupResp, err := this.RPC().NodeGroupRPC().FindEnabledNodeGroup(this.AdminContext(), &pb.FindEnabledNodeGroupRequest{NodeGroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	group := groupResp.NodeGroup
	if group == nil {
		this.NotFound("nodeGroup", params.GroupId)
		return
	}

	this.Data["group"] = maps.Map{
		"id":   group.Id,
		"name": group.Name,
	}

	this.Success()
}
