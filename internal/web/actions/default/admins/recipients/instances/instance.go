package instances

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/monitorconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type InstanceAction struct {
	actionutils.ParentAction
}

func (this *InstanceAction) Init() {
	this.Nav("", "", "instance")
}

func (this *InstanceAction) RunGet(params struct {
	InstanceId int64
}) {
	instanceResp, err := this.RPC().MessageMediaInstanceRPC().FindEnabledMessageMediaInstance(this.AdminContext(), &pb.FindEnabledMessageMediaInstanceRequest{MessageMediaInstanceId: params.InstanceId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	instance := instanceResp.MessageMediaInstance
	if instance == nil || instance.MessageMedia == nil {
		this.NotFound("messageMediaInstance", params.InstanceId)
		return
	}

	mediaParams := maps.Map{}
	if len(instance.ParamsJSON) > 0 {
		err = json.Unmarshal(instance.ParamsJSON, &mediaParams)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// 频率
	var rateConfig = &monitorconfigs.RateConfig{}
	if len(instance.RateJSON) > 0 {
		err = json.Unmarshal(instance.RateJSON, rateConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["instance"] = maps.Map{
		"id":   instance.Id,
		"name": instance.Name,
		"isOn": instance.IsOn,
		"media": maps.Map{
			"type": instance.MessageMedia.Type,
			"name": instance.MessageMedia.Name,
		},
		"description": instance.Description,
		"params":      mediaParams,
		"rate":        rateConfig,
		"hashLife":    instance.HashLife,
	}

	this.Show()
}
