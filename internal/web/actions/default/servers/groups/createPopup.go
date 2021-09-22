package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")
	createResp, err := this.RPC().ServerGroupRPC().CreateServerGroup(this.AdminContext(), &pb.CreateServerGroupRequest{
		Name: params.Name,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["group"] = maps.Map{
		"id":   createResp.ServerGroupId,
		"name": params.Name,
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "创建代理服务分组 %d", createResp.ServerGroupId)

	this.Success()
}
