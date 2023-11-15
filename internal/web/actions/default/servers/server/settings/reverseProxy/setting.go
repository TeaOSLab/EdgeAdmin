package reverseProxy

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
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
	var reverseProxyRef = &serverconfigs.ReverseProxyRef{}
	err = json.Unmarshal(reverseProxyResp.ReverseProxyRefJSON, reverseProxyRef)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var reverseProxy = serverconfigs.NewReverseProxyConfig()
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
	defer this.CreateLogInfo(codes.ServerReverseProxy_LogUpdateServerReverseProxySettings, params.ServerId)

	var reverseProxyConfig = serverconfigs.NewReverseProxyConfig()
	err := json.Unmarshal(params.ReverseProxyJSON, reverseProxyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	err = reverseProxyConfig.Init(context.TODO())
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
	}

	if reverseProxyConfig.ConnTimeout == nil {
		reverseProxyConfig.ConnTimeout = &shared.TimeDuration{Count: 0, Unit: "second"}
	}
	connTimeoutJSON, err := json.Marshal(reverseProxyConfig.ConnTimeout)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if reverseProxyConfig.ReadTimeout == nil {
		reverseProxyConfig.ReadTimeout = &shared.TimeDuration{Count: 0, Unit: "second"}
	}
	readTimeoutJSON, err := json.Marshal(reverseProxyConfig.ReadTimeout)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if reverseProxyConfig.IdleTimeout == nil {
		reverseProxyConfig.IdleTimeout = &shared.TimeDuration{Count: 0, Unit: "second"}
	}
	idleTimeoutJSON, err := json.Marshal(reverseProxyConfig.IdleTimeout)
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

	// PROXY Protocol
	var proxyProtocolJSON = []byte{}
	if reverseProxyConfig.ProxyProtocol != nil {
		proxyProtocolJSON, err = json.Marshal(reverseProxyConfig.ProxyProtocol)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// 设置反向代理相关信息
	_, err = this.RPC().ReverseProxyRPC().UpdateReverseProxy(this.AdminContext(), &pb.UpdateReverseProxyRequest{
		ReverseProxyId:           reverseProxyConfig.Id,
		RequestHostType:          types.Int32(reverseProxyConfig.RequestHostType),
		RequestHost:              reverseProxyConfig.RequestHost,
		RequestURI:               reverseProxyConfig.RequestURI,
		StripPrefix:              reverseProxyConfig.StripPrefix,
		AutoFlush:                reverseProxyConfig.AutoFlush,
		AddHeaders:               reverseProxyConfig.AddHeaders,
		ConnTimeoutJSON:          connTimeoutJSON,
		ReadTimeoutJSON:          readTimeoutJSON,
		IdleTimeoutJSON:          idleTimeoutJSON,
		MaxConns:                 types.Int32(reverseProxyConfig.MaxConns),
		MaxIdleConns:             types.Int32(reverseProxyConfig.MaxIdleConns),
		ProxyProtocolJSON:        proxyProtocolJSON,
		FollowRedirects:          reverseProxyConfig.FollowRedirects,
		RequestHostExcludingPort: reverseProxyConfig.RequestHostExcludingPort,
		Retry50X:                 reverseProxyConfig.Retry50X,
		Retry40X:                 reverseProxyConfig.Retry40X,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
