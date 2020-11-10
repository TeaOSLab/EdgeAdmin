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
	ClusterId int64
	Name      string

	Must *actions.Must
}) {
	if params.ClusterId <= 0 {
		this.Fail("请选择集群")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")
	createResp, err := this.RPC().NodeGroupRPC().CreateNodeGroup(this.AdminContext(), &pb.CreateNodeGroupRequest{
		ClusterId: params.ClusterId,
		Name:      params.Name,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["group"] = maps.Map{
		"id":   createResp.GroupId,
		"name": params.Name,
	}

	// 创建日志
	this.CreateLog(oplogs.LevelInfo, "创建集群分组", createResp.GroupId)

	this.Success()
}
