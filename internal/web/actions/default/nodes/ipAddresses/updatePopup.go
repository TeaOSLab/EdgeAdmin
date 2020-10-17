package ipAddresses

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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

	Must *actions.Must
}) {
	// TODO 严格校验IP地址

	params.Must.
		Field("ip", params.IP).
		Require("请输入IP地址")

	ip := net.ParseIP(params.IP)
	if len(ip) == 0 {
		this.Fail("请输入正确的IP")
	}

	this.Data["ipAddress"] = maps.Map{
		"name":      params.Name,
		"ip":        params.IP,
		"id":        params.AddressId,
		"canAccess": params.CanAccess,
	}

	this.Success()
}
