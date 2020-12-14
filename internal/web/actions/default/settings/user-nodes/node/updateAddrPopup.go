package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"net"
)

type UpdateAddrPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateAddrPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateAddrPopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *UpdateAddrPopupAction) RunPost(params struct {
	Protocol string
	Addr     string
	Must     *actions.Must
}) {
	params.Must.
		Field("addr", params.Addr).
		Require("请输入访问地址")
	host, port, err := net.SplitHostPort(params.Addr)
	if err != nil {
		this.FailField("addr", "错误的访问地址")
	}

	addrConfig := &serverconfigs.NetworkAddressConfig{
		Protocol:  serverconfigs.Protocol(params.Protocol),
		Host:      host,
		PortRange: port,
	}
	this.Data["addr"] = addrConfig
	this.Success()
}
