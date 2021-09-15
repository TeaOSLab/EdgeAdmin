package ipAddresses

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
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
	AddressId         int64
	SupportThresholds bool
}) {
	this.Data["supportThresholds"] = params.SupportThresholds

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	AddressId      int64
	IP             string `alias:"ip"`
	Name           string
	CanAccess      bool
	IsOn           bool
	IsUp           bool
	ThresholdsJSON []byte

	Must *actions.Must
}) {
	params.Must.
		Field("ip", params.IP).
		Require("请输入IP地址")

	// 获取IP地址信息
	var isUp = params.IsUp
	if params.AddressId > 0 {
		addressResp, err := this.RPC().NodeIPAddressRPC().FindEnabledNodeIPAddress(this.AdminContext(), &pb.FindEnabledNodeIPAddressRequest{NodeIPAddressId: params.AddressId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var address = addressResp.NodeIPAddress
		if address == nil {
			this.Fail("找不到要修改的地址")
		}
	}

	ip := net.ParseIP(params.IP)
	if len(ip) == 0 {
		this.Fail("请输入正确的IP")
	}

	var thresholds = []*nodeconfigs.IPAddressThresholdConfig{}
	if teaconst.IsPlus && len(params.ThresholdsJSON) > 0 {
		_ = json.Unmarshal(params.ThresholdsJSON, &thresholds)
	}

	this.Data["ipAddress"] = maps.Map{
		"name":       params.Name,
		"ip":         params.IP,
		"id":         params.AddressId,
		"canAccess":  params.CanAccess,
		"isOn":       params.IsOn,
		"isUp":       isUp,
		"thresholds": thresholds,
	}

	this.Success()
}
