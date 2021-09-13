package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

// StatusAction 节点状态
type StatusAction struct {
	actionutils.ParentAction
}

func (this *StatusAction) RunPost(params struct {
	NodeId int64
}) {
	// 节点
	nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	node := nodeResp.Node
	if node == nil {
		this.WriteString("找不到要操作的节点")
		return
	}

	// 安装信息
	if node.InstallStatus != nil {
		this.Data["installStatus"] = maps.Map{
			"isRunning":  node.InstallStatus.IsRunning,
			"isFinished": node.InstallStatus.IsFinished,
			"isOk":       node.InstallStatus.IsOk,
			"updatedAt":  node.InstallStatus.UpdatedAt,
			"error":      node.InstallStatus.Error,
			"errorCode":  node.InstallStatus.ErrorCode,
		}
	} else {
		this.Data["installStatus"] = nil
	}

	this.Data["isInstalled"] = node.IsInstalled

	this.Success()
}
