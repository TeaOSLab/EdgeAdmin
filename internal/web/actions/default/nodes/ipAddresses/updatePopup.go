package ipAddresses

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	AddressId int64
}) {
	addressResp, err := this.RPC().NodeIPAddressRPC().FindEnabledNodeIPAddress(this.AdminContext(), &pb.FindEnabledNodeIPAddressRequest{AddressId: params.AddressId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	address := addressResp.IpAddress
	if address == nil {
		this.WriteString("找不到要修改的IP地址")
		return
	}

	this.Data["address"] = maps.Map{
		"id":   address.Id,
		"name": address.Name,
		"ip":   address.Ip,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	AddressId int64
	IP        string `alias:"ip"`
	Name      string

	Must *actions.Must
}) {
	// TODO 严格校验IP地址

	params.Must.
		Field("ip", params.IP).
		Require("请输入IP地址")

	_, err := this.RPC().NodeIPAddressRPC().UpdateNodeIPAddress(this.AdminContext(), &pb.UpdateNodeIPAddressRequest{
		AddressId: params.AddressId,
		Name:      params.Name,
		Ip:        params.IP,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["ipAddress"] = maps.Map{
		"name": params.Name,
		"ip":   params.IP,
		"id":   params.AddressId,
	}

	this.Success()
}
