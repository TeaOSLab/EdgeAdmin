package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"net"
	"net/url"
	"regexp"
	"strings"
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

	// 兼容URL
	if regexp.MustCompile(`^(?i)(http|https)://`).MatchString(params.Addr) {
		u, err := url.Parse(params.Addr)
		if err != nil {
			this.FailField("addr", "错误的访问地址，不需要添加http://或https://")
		}
		params.Addr = u.Host
	}

	// 自动添加端口
	if !strings.Contains(params.Addr, ":") {
		switch params.Protocol {
		case "http":
			params.Addr += ":80"
		case "https":
			params.Addr += ":443"
		}
	}

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
