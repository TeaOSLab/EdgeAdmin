package settings

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type HealthAction struct {
	actionutils.ParentAction
}

func (this *HealthAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("health")
}

func (this *HealthAction) RunGet(params struct {
	ClusterId int64
}) {
	configResp, err := this.RPC().NodeClusterRPC().FindNodeClusterHealthCheckConfig(this.AdminContext(), &pb.FindNodeClusterHealthCheckConfigRequest{ClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var config *serverconfigs.HealthCheckConfig = nil
	if len(configResp.HealthCheckConfig) > 0 {
		config = &serverconfigs.HealthCheckConfig{}
		err = json.Unmarshal(configResp.HealthCheckConfig, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["healthCheckConfig"] = config
	this.Show()
}

func (this *HealthAction) RunPost(params struct {
	ClusterId       int64
	HealthCheckJSON []byte
	Must            *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改集群健康检查设置 %d", params.ClusterId)

	config := &serverconfigs.HealthCheckConfig{}
	err := json.Unmarshal(params.HealthCheckJSON, config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().NodeClusterRPC().UpdateNodeClusterHealthCheck(this.AdminContext(), &pb.UpdateNodeClusterHealthCheckRequest{
		ClusterId:       params.ClusterId,
		HealthCheckJSON: params.HealthCheckJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
