package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"net/url"
	"regexp"
	"strings"
)

type AddOriginPopupAction struct {
	actionutils.ParentAction
}

func (this *AddOriginPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *AddOriginPopupAction) RunGet(params struct {
	ServerType string
}) {
	this.Data["serverType"] = params.ServerType

	this.Show()
}

func (this *AddOriginPopupAction) RunPost(params struct {
	Protocol string
	Addr     string

	Must *actions.Must
}) {
	params.Must.
		Field("addr", params.Addr).
		Require("请输入源站地址")

	addr := params.Addr

	// 是否是完整的地址
	if (params.Protocol == "http" || params.Protocol == "https") && regexp.MustCompile(`^(http|https)://`).MatchString(addr) {
		u, err := url.Parse(addr)
		if err == nil {
			addr = u.Host
		}
	}

	addr = regexp.MustCompile(`\s+`).ReplaceAllString(addr, "")
	portIndex := strings.LastIndex(addr, ":")
	if portIndex < 0 {
		if params.Protocol == "http" {
			addr += ":80"
		} else if params.Protocol == "https" {
			addr += ":443"
		} else {
			this.Fail("地址中需要带有端口")
		}
		portIndex = strings.LastIndex(addr, ":")
	}
	host := addr[:portIndex]
	port := addr[portIndex+1:]

	resp, err := this.RPC().OriginRPC().CreateOrigin(this.AdminContext(), &pb.CreateOriginRequest{
		Name: "",
		Addr: &pb.NetworkAddress{
			Protocol:  params.Protocol,
			Host:      host,
			PortRange: port,
		},
		Description: "",
		Weight:      10,
		IsOn:        true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	origin := &serverconfigs.OriginConfig{
		Id:   resp.OriginId,
		IsOn: true,
		Addr: &serverconfigs.NetworkAddressConfig{
			Protocol:  serverconfigs.Protocol(params.Protocol),
			Host:      host,
			PortRange: port,
		},
	}

	this.Data["origin"] = origin

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "创建源站 %d", resp.OriginId)

	this.Success()
}
