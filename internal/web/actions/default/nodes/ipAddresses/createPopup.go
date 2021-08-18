package ipAddresses

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
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
	IP             string `alias:"ip"`
	CanAccess      bool
	Name           string
	ThresholdsJSON []byte

	Must *actions.Must
}) {
	ip := net.ParseIP(params.IP)
	if len(ip) == 0 {
		this.Fail("请输入正确的IP")
	}

	params.Must.
		Field("ip", params.IP).
		Require("请输入IP地址")

	var thresholds = []*nodeconfigs.NodeValueThresholdConfig{}
	if teaconst.IsPlus && len(params.ThresholdsJSON) > 0 {
		_ = json.Unmarshal(params.ThresholdsJSON, &thresholds)
	}

	this.Data["ipAddress"] = maps.Map{
		"name":       params.Name,
		"canAccess":  params.CanAccess,
		"ip":         params.IP,
		"id":         0,
		"isOn":       true,
		"isUp":       true,
		"thresholds": thresholds,
	}
	this.Success()
}
