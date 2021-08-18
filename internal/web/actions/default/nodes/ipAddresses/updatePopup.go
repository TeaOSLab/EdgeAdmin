package ipAddresses

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net"
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
	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	AddressId int64
	IP        string `alias:"ip"`
	Name      string
	CanAccess bool
	IsOn      bool

	Must *actions.Must
}) {
	params.Must.
		Field("ip", params.IP).
		Require("请输入IP地址")

	// 获取IP地址信息
	addressResp, err := this.RPC().NodeIPAddressRPC().FindEnabledNodeIPAddress(this.AdminContext(), &pb.FindEnabledNodeIPAddressRequest{AddressId: params.AddressId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var address = addressResp.IpAddress
	if address == nil {
		this.Fail("找不到要修改的地址")
	}

	ip := net.ParseIP(params.IP)
	if len(ip) == 0 {
		this.Fail("请输入正确的IP")
	}

	this.Data["ipAddress"] = maps.Map{
		"name":      params.Name,
		"ip":        params.IP,
		"id":        params.AddressId,
		"canAccess": params.CanAccess,
		"isOn":      params.IsOn,
		"isUp":      address.IsUp,
	}

	this.Success()
}
