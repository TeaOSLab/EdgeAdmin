package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type UpgradeStatusAction struct {
	actionutils.ParentAction
}

func (this *UpgradeStatusAction) RunPost(params struct {
	NodeId int64
}) {
	resp, err := this.RPC().NodeRPC().FindNodeInstallStatus(this.AdminContext(), &pb.FindNodeInstallStatusRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if resp.InstallStatus == nil {
		this.Data["status"] = nil
		this.Success()
	}

	this.Data["status"] = maps.Map{
		"isRunning":  resp.InstallStatus.IsRunning,
		"isFinished": resp.InstallStatus.IsFinished,
		"isOk":       resp.InstallStatus.IsOk,
		"error":      resp.InstallStatus.Error,
		"errorCode":  resp.InstallStatus.ErrorCode,
	}

	this.Success()
}
