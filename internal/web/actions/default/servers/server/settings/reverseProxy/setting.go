package reverseProxy

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type SettingAction struct {
	actionutils.ParentAction
}

func (this *SettingAction) Init() {
	this.FirstMenu("setting")
}

func (this *SettingAction) RunGet(params struct {
	ServerId int64
}) {
	reverseProxyResp, err := this.RPC().ServerRPC().FindAndInitServerReverseProxyConfig(this.AdminContext(), &pb.FindAndInitServerReverseProxyConfigRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	reverseProxyRef := &serverconfigs.ReverseProxyRef{}
	err = json.Unmarshal(reverseProxyResp.ReverseProxyRefJSON, reverseProxyRef)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	reverseProxy := &serverconfigs.ReverseProxyConfig{}
	err = json.Unmarshal(reverseProxyResp.ReverseProxyJSON, reverseProxy)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["reverseProxyRef"] = reverseProxyRef
	this.Data["reverseProxyConfig"] = reverseProxy

	this.Show()
}

func (this *SettingAction) RunPost(params struct {
	ServerId            int64
	ReverseProxyRefJSON []byte
	ReverseProxyJSON    []byte

	Must *actions.Must
}) {
	// TODO 校验配置

	reverseProxyConfig := &serverconfigs.ReverseProxyConfig{}
	err := json.Unmarshal(params.ReverseProxyJSON, reverseProxyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 设置是否启用
	_, err = this.RPC().ServerRPC().UpdateServerReverseProxy(this.AdminContext(), &pb.UpdateServerReverseProxyRequest{
		ServerId:         params.ServerId,
		ReverseProxyJSON: params.ReverseProxyRefJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 设置反向代理相关信息
	_, err = this.RPC().ReverseProxyRPC().UpdateReverseProxy(this.AdminContext(), &pb.UpdateReverseProxyRequest{
		ReverseProxyId: reverseProxyConfig.Id,
		RequestHost:    reverseProxyConfig.RequestHost,
		RequestURI:     reverseProxyConfig.RequestURI,
		StripPrefix:    reverseProxyConfig.StripPrefix,
		AutoFlush:      reverseProxyConfig.AutoFlush,
	})

	this.Success()
}
