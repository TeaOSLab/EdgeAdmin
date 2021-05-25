package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	GroupId int64
}) {
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

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	GroupId int64
	Name    string

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改代理服务分组 %d", params.GroupId)

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
