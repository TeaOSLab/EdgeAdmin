package servers

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
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

	addr := regexp.MustCompile(`\s+`).ReplaceAllString(params.Addr, "")
	portIndex := strings.LastIndex(params.Addr, ":")
	if portIndex < 0 {
		this.Fail("地址中需要带有端口")
	}
	host := addr[:portIndex]
	port := addr[portIndex+1:]

	resp, err := this.RPC().OriginServerRPC().CreateOriginServer(this.AdminContext(), &pb.CreateOriginServerRequest{
		Name: "",
		Addr: &pb.NetworkAddress{
			Protocol:  params.Protocol,
			Host:      host,
			PortRange: port,
		},
		Description: "",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	origin := &serverconfigs.OriginServerConfig{
		Id:   resp.OriginId,
		IsOn: true,
		Addr: &serverconfigs.NetworkAddressConfig{
			Protocol:  params.Protocol,
			Host:      host,
			PortRange: port,
		},
	}

	this.Data["origin"] = origin
	this.Success()
}
