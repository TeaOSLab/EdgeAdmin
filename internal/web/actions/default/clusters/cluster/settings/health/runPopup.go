package health

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type RunPopupAction struct {
	actionutils.ParentAction
}

func (this *RunPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *RunPopupAction) RunGet(params struct {
	ClusterId int64
}) {
	// 检查是否已部署服务
	countServersResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithNodeClusterId(this.AdminContext(), &pb.CountAllEnabledServersWithNodeClusterIdRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasServers"] = countServersResp.Count > 0

	this.Show()
}

func (this *RunPopupAction) RunPost(params struct {
	ClusterId int64

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "执行集群健康检查设置 %d", params.ClusterId)

	resp, err := this.RPC().NodeClusterRPC().ExecuteNodeClusterHealthCheck(this.AdminContext(), &pb.ExecuteNodeClusterHealthCheckRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.Fail(err.Error())
	}

	if resp.Results == nil {
		resp.Results = []*pb.ExecuteNodeClusterHealthCheckResponse_Result{}
	}
	this.Data["results"] = resp.Results
	this.Success()
}
