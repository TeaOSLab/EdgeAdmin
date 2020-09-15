package reverseProxy

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/schedulingconfigs"
)

type SchedulingAction struct {
	actionutils.ParentAction
}

func (this *SchedulingAction) Init() {
	this.FirstMenu("scheduling")
}

func (this *SchedulingAction) RunGet(params struct {
	ServerId int64
}) {
	server, _, isOk := serverutils.FindServer(&this.ParentAction, params.ServerId)
	if !isOk {
		return
	}

	if server.ReverseProxyId <= 0 {
		// TODO 在界面上提示用户未开通，并提供开通按钮，用户点击后开通
		this.WriteString("此服务尚未开通反向代理功能")
		return
	}
	this.Data["reverseProxyId"] = server.ReverseProxyId

	reverseProxyResp, err := this.RPC().ReverseProxyRPC().FindEnabledReverseProxyConfig(this.AdminContext(), &pb.FindEnabledReverseProxyConfigRequest{
		ReverseProxyId: server.ReverseProxyId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	reverseProxy := &serverconfigs.ReverseProxyConfig{}
	err = json.Unmarshal(reverseProxyResp.Config, reverseProxy)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	schedulingCode := reverseProxy.FindSchedulingConfig().Code
	schedulingMap := schedulingconfigs.FindSchedulingType(schedulingCode)
	if schedulingMap == nil {
		this.ErrorPage(errors.New("invalid scheduling code '" + schedulingCode + "'"))
		return
	}
	this.Data["scheduling"] = schedulingMap

	this.Show()
}
