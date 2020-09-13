package ipAddresses

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	IP   string `alias:"ip"`
	Name string

	Must *actions.Must
}) {
	// TODO 严格校验IP地址

	params.Must.
		Field("ip", params.IP).
		Require("请输入IP地址")

	resp, err := this.RPC().NodeIPAddressRPC().CreateNodeIPAddress(this.AdminContext(), &pb.CreateNodeIPAddressRequest{
		NodeId: 0,
		Name:   params.Name,
		Ip:     params.IP,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["ipAddress"] = maps.Map{
		"name": params.Name,
		"ip":   params.IP,
		"id":   resp.AddressId,
	}
	this.Success()
}
