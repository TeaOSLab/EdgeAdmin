package ipAddresses

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/ipAddresses/ipaddressutils"
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
	NodeId            int64
	SupportThresholds bool
}) {
	// 专属集群
	clusterMaps, err := ipaddressutils.FindNodeClusterMapsWithNodeId(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["clusters"] = clusterMaps

	// 阈值
	this.Data["supportThresholds"] = params.SupportThresholds

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	IP             string `alias:"ip"`
	CanAccess      bool
	Name           string
	IsUp           bool
	ThresholdsJSON []byte
	ClusterIds     []int64

	Must *actions.Must
}) {
	params.Must.
		Field("ip", params.IP).
		Require("请输入IP地址")

	result, err := utils.ExtractIP(params.IP)
	if err != nil {
		this.Fail("IP格式错误'" + params.IP + "'")
	}

	for _, ip := range result {
		if len(net.ParseIP(ip)) == 0 {
			this.FailField("ip", "请输入正确的IP")
		}
	}

	// 阈值设置
	var thresholds = []*nodeconfigs.IPAddressThresholdConfig{}
	if teaconst.IsPlus && len(params.ThresholdsJSON) > 0 {
		_ = json.Unmarshal(params.ThresholdsJSON, &thresholds)
	}

	// 专属集群
	// 目前只考虑CDN边缘集群
	clusterMaps, err := ipaddressutils.FindNodeClusterMaps(this.Parent(), params.ClusterIds)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["ipAddress"] = maps.Map{
		"name":       params.Name,
		"canAccess":  params.CanAccess,
		"ip":         params.IP,
		"id":         0,
		"isOn":       true,
		"isUp":       params.IsUp,
		"thresholds": thresholds,
		"clusters":   clusterMaps,
	}
	this.Success()
}
