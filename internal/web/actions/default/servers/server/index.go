package server

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"strconv"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "index", "index")
	this.SecondMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// TODO 等看板实现后，需要跳转到看板

	// TCP & UDP跳转到设置
	serverTypeResp, err := this.RPC().ServerRPC().FindEnabledServerType(this.AdminContext(), &pb.FindEnabledServerTypeRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	serverType := serverTypeResp.Type
	if serverType == serverconfigs.ServerTypeTCPProxy || serverType == serverconfigs.ServerTypeUDPProxy {
		this.RedirectURL("/servers/server/settings?serverId=" + strconv.FormatInt(params.ServerId, 10))
		return
	}

	// HTTP跳转到访问日志
	this.RedirectURL("/servers/server/log?serverId=" + strconv.FormatInt(params.ServerId, 10))
}
