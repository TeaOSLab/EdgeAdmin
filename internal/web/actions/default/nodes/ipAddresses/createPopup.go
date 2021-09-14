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

func (this *CreatePopupAction) RunGet(params struct {
	SupportThresholds bool
}) {
	this.Data["supportThresholds"] = params.SupportThresholds

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	IP             string `alias:"ip"`
	CanAccess      bool
	Name           string
	IsUp           bool
	ThresholdsJSON []byte

	Must *actions.Must
}) {
	params.Must.
		Field("ip", params.IP).
		Require("请输入IP地址")

	ip := net.ParseIP(params.IP)
	if len(ip) == 0 {
		this.FailField("ip", "请输入正确的IP")
	}

	// 阈值设置
	var thresholds = []*nodeconfigs.IPAddressThresholdConfig{}
	if teaconst.IsPlus && len(params.ThresholdsJSON) > 0 {
		_ = json.Unmarshal(params.ThresholdsJSON, &thresholds)
	}

	this.Data["ipAddress"] = maps.Map{
		"name":       params.Name,
		"canAccess":  params.CanAccess,
		"ip":         params.IP,
		"id":         0,
		"isOn":       true,
		"isUp":       params.IsUp,
		"thresholds": thresholds,
	}
	this.Success()
}
