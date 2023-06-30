package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
		NodeClusterId: params.ClusterId,
		Name:          params.Name,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["group"] = maps.Map{
		"id":   createResp.NodeGroupId,
		"name": params.Name,
	}

	// 创建日志
	defer this.CreateLogInfo(codes.NodeGroup_LogCreateNodeGroup, createResp.NodeGroupId)

	this.Success()
}
