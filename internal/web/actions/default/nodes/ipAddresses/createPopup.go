package ipAddresses

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net"
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
	IP        string `alias:"ip"`
	CanAccess bool
	Name      string

	Must *actions.Must
}) {
	// TODO 严格校验IP地址

	ip := net.ParseIP(params.IP)
	if len(ip) == 0 {
		this.Fail("请输入正确的IP")
	}

	params.Must.
		Field("ip", params.IP).
		Require("请输入IP地址")

	this.Data["ipAddress"] = maps.Map{
		"name":      params.Name,
		"canAccess": params.CanAccess,
		"ip":        params.IP,
		"id":        0,
		"isOn":      true,
		"isUp":      true,
	}
	this.Success()
}
