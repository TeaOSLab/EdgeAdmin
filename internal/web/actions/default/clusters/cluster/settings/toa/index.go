package toa

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("toa")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	toaResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterTOA(this.AdminContext(), &pb.FindEnabledNodeClusterTOARequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if len(toaResp.ToaJSON) == 0 {
		this.Data["toa"] = nodeconfigs.DefaultTOAConfig()
	} else {
		config := &nodeconfigs.TOAConfig{}
		err = json.Unmarshal(toaResp.ToaJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.Data["toa"] = config
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId  int64
	IsOn       bool
	AutoSetup  bool
	OptionType uint8
	MinQueueId uint8
	MaxQueueId uint8

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改集群 %d 的TOA设置", params.ClusterId)

	config := &nodeconfigs.TOAConfig{
		IsOn:       params.IsOn,
		Debug:      false, // 暂时不允许打开调试
		OptionType: params.OptionType,
		MinQueueId: params.MinQueueId,
		MaxQueueId: params.MaxQueueId,
		AutoSetup:  params.AutoSetup,
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().NodeClusterRPC().UpdateNodeClusterTOA(this.AdminContext(), &pb.UpdateNodeClusterTOARequest{
		NodeClusterId: params.ClusterId,
		ToaJSON:       configJSON,
	})

	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
