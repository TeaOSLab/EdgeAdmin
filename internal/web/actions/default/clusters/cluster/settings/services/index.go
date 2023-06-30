package services

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "setting")
	this.SecondMenu("service")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	serviceParamsResp, err := this.RPC().NodeClusterRPC().FindNodeClusterSystemService(this.AdminContext(), &pb.FindNodeClusterSystemServiceRequest{
		NodeClusterId: params.ClusterId,
		Type:          nodeconfigs.SystemServiceTypeSystemd,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	paramsJSON := serviceParamsResp.ParamsJSON
	if len(paramsJSON) == 0 {
		this.Data["systemdIsOn"] = false
	} else {
		config := &nodeconfigs.SystemdServiceConfig{}
		err = json.Unmarshal(paramsJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.Data["systemdIsOn"] = config.IsOn
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId   int64
	SystemdIsOn bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.NodeSystemd_LogUpdateClusterSystemdSettings, params.ClusterId)

	serviceParams := &nodeconfigs.SystemdServiceConfig{
		IsOn: params.SystemdIsOn,
	}
	serviceParamsJSON, err := json.Marshal(serviceParams)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().NodeClusterRPC().UpdateNodeClusterSystemService(this.AdminContext(), &pb.UpdateNodeClusterSystemServiceRequest{
		NodeClusterId: params.ClusterId,
		Type:          nodeconfigs.SystemServiceTypeSystemd,
		ParamsJSON:    serviceParamsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
