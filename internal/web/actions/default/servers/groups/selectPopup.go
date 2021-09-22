package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectPopupAction) RunGet(params struct {
	SelectedGroupIds string
}) {
	groupsResp, err := this.RPC().ServerGroupRPC().FindAllEnabledServerGroups(this.AdminContext(), &pb.FindAllEnabledServerGroupsRequest{})
	if err != nil {
		this.ErrorPage(err)
	}

	selectedGroupIds := []int64{}
	if len(params.SelectedGroupIds) > 0 {
		for _, v := range strings.Split(params.SelectedGroupIds, ",") {
			selectedGroupIds = append(selectedGroupIds, types.Int64(v))
		}
	}

	groupMaps := []maps.Map{}
	for _, group := range groupsResp.ServerGroups {
		// 已经选过的就跳过
		if lists.ContainsInt64(selectedGroupIds, group.Id) {
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

func (this *SelectPopupAction) RunPost(params struct {
	GroupId int64

	Must *actions.Must
}) {
	if params.GroupId <= 0 {
		this.Fail("请选择要使用的分组")
	}

	groupResp, err := this.RPC().ServerGroupRPC().FindEnabledServerGroup(this.AdminContext(), &pb.FindEnabledServerGroupRequest{ServerGroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	group := groupResp.ServerGroup
	if group == nil {
		this.NotFound("serverGroup", params.GroupId)
		return
	}

	this.Data["group"] = maps.Map{
		"id":   group.Id,
		"name": group.Name,
	}

	this.Success()
}
