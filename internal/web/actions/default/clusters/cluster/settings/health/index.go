package health

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("health")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	configResp, err := this.RPC().NodeClusterRPC().FindNodeClusterHealthCheckConfig(this.AdminContext(), &pb.FindNodeClusterHealthCheckConfigRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var config *serverconfigs.HealthCheckConfig = nil
	if len(configResp.HealthCheckJSON) > 0 {
		config = &serverconfigs.HealthCheckConfig{}
		err = json.Unmarshal(configResp.HealthCheckJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["healthCheckConfig"] = config
	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId       int64
	HealthCheckJSON []byte
	Must            *actions.Must
}) {
	// 创建日志
	defer this.CreateLogInfo(codes.NodeCluster_LogUpdateClusterHealthCheck, params.ClusterId)

	config := &serverconfigs.HealthCheckConfig{}
	err := json.Unmarshal(params.HealthCheckJSON, config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().NodeClusterRPC().UpdateNodeClusterHealthCheck(this.AdminContext(), &pb.UpdateNodeClusterHealthCheckRequest{
		NodeClusterId:   params.ClusterId,
		HealthCheckJSON: params.HealthCheckJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
