package group

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/groups/group/servergrouputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "group.update")
}

func (this *UpdateAction) RunGet(params struct {
	GroupId int64
}) {
	group, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["group"] = maps.Map{
		"id":   group.Id,
		"name": group.Name,
	}

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	GroupId int64
	Name    string

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLogInfo(codes.ServerGroup_LogUpdateServerGroup, params.GroupId)

	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")
	_, err := this.RPC().ServerGroupRPC().UpdateServerGroup(this.AdminContext(), &pb.UpdateServerGroupRequest{
		ServerGroupId: params.GroupId,
		Name:          params.Name,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
