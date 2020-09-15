package reverseProxy

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/schedulingconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

// 修改调度算法
type UpdateSchedulingPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateSchedulingPopupAction) Init() {
}

func (this *UpdateSchedulingPopupAction) RunGet(params struct {
	Type           string
	ServerId       int64
	ReverseProxyId int64
}) {
	this.Data["dataType"] = params.Type
	this.Data["serverId"] = params.ServerId
	this.Data["reverseProxyId"] = params.ReverseProxyId

	_, serverConfig, isOk := serverutils.FindServer(&this.ParentAction, params.ServerId)
	if !isOk {
		return
	}

	reverseProxyResp, err := this.RPC().ReverseProxyRPC().FindEnabledReverseProxyConfig(this.AdminContext(), &pb.FindEnabledReverseProxyConfigRequest{
		ReverseProxyId: params.ReverseProxyId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	configData := reverseProxyResp.Config

	reverseProxyConfig := &serverconfigs.ReverseProxyConfig{}
	err = json.Unmarshal(configData, reverseProxyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	schedulingObject := &serverconfigs.SchedulingConfig{
		Code:    "random",
		Options: nil,
	}
	if reverseProxyConfig.Scheduling != nil {
		schedulingObject = reverseProxyConfig.Scheduling
	}
	this.Data["scheduling"] = schedulingObject

	// 调度类型
	schedulingTypes := []maps.Map{}
	for _, m := range schedulingconfigs.AllSchedulingTypes() {
		networks, ok := m["networks"]
		if !ok {
			continue
		}
		if !types.IsSlice(networks) {
			continue
		}
		if (serverConfig.IsHTTP() && lists.Contains(networks, "http")) ||
			(serverConfig.IsTCP() && lists.Contains(networks, "tcp")) ||
			(serverConfig.IsUDP() && lists.Contains(networks, "udp")) ||
			(serverConfig.IsUnix() && lists.Contains(networks, "unix")) {
			schedulingTypes = append(schedulingTypes, m)
		}
	}
	this.Data["schedulingTypes"] = schedulingTypes

	this.Show()
}

func (this *UpdateSchedulingPopupAction) RunPost(params struct {
	ServerId       int64
	ReverseProxyId int64

	Type        string
	HashKey     string
	StickyType  string
	StickyParam string

	Must *actions.Must
}) {
	reverseProxyResp, err := this.RPC().ReverseProxyRPC().FindEnabledReverseProxyConfig(this.AdminContext(), &pb.FindEnabledReverseProxyConfigRequest{ReverseProxyId: params.ReverseProxyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	configData := reverseProxyResp.Config
	reverseProxy := &serverconfigs.ReverseProxyConfig{}
	err = json.Unmarshal(configData, reverseProxy)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if reverseProxy.Scheduling == nil {
		reverseProxy.FindSchedulingConfig()
	}

	options := maps.Map{}
	if params.Type == "hash" {
		params.Must.
			Field("hashKey", params.HashKey).
			Require("请输入Key")

		options["key"] = params.HashKey
	} else if params.Type == "sticky" {
		params.Must.
			Field("stickyType", params.StickyType).
			Require("请选择参数类型").
			Field("stickyParam", params.StickyParam).
			Require("请输入参数名").
			Match("^[a-zA-Z0-9]+$", "参数名只能是英文字母和数字的组合").
			MaxCharacters(50, "参数名长度不能超过50位")

		options["type"] = params.StickyType
		options["param"] = params.StickyParam
	}

	if schedulingconfigs.FindSchedulingType(params.Type) == nil {
		this.Fail("不支持此种算法")
	}

	reverseProxy.Scheduling.Code = params.Type
	reverseProxy.Scheduling.Options = options

	schedulingData, err := json.Marshal(reverseProxy.Scheduling)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ReverseProxyRPC().UpdateReverseProxyScheduling(this.AdminContext(), &pb.UpdateReverseProxySchedulingRequest{
		ReverseProxyId: params.ReverseProxyId,
		SchedulingJSON: schedulingData,
	})
	if err != nil {
		this.ErrorPage(err)
	}

	this.Success()
}
