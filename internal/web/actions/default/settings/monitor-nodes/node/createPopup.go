package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "node", "create")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name        string
	Description string
	IsOn        bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入监控节点名称")

	createResp, err := this.RPC().MonitorNodeRPC().CreateMonitorNode(this.AdminContext(), &pb.CreateMonitorNodeRequest{
		Name:        params.Name,
		Description: params.Description,
		IsOn:        params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "创建监控节点 %d", createResp.NodeId)

	this.Success()
}
