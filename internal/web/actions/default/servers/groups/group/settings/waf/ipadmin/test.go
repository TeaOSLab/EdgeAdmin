package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type TestAction struct {
	actionutils.ParentAction
}

func (this *TestAction) Init() {
	this.Nav("", "setting", "test")
	this.SecondMenu("waf")
}

func (this *TestAction) RunGet(params struct {
	ServerId         int64
	FirewallPolicyId int64
}) {
	this.Data["featureIsOn"] = true
	this.Data["firewallPolicyId"] = params.FirewallPolicyId
	this.Data["subMenuItem"] = "province"

	// WAF是否启用
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["wafIsOn"] = webConfig.FirewallRef != nil && webConfig.FirewallRef.IsOn

	this.Show()
}

func (this *TestAction) RunPost(params struct {
	FirewallPolicyId int64
	Ip               string

	Must *actions.Must
}) {
	resp, err := this.RPC().HTTPFirewallPolicyRPC().CheckHTTPFirewallPolicyIPStatus(this.AdminContext(), &pb.CheckHTTPFirewallPolicyIPStatusRequest{
		HttpFirewallPolicyId: params.FirewallPolicyId,
		Ip:                   params.Ip,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	resultMap := maps.Map{
		"isDone":    true,
		"isFound":   resp.IsFound,
		"isOk":      resp.IsOk,
		"error":     resp.Error,
		"isAllowed": resp.IsAllowed,
	}

	if resp.IpList != nil {
		resultMap["list"] = maps.Map{
			"id":   resp.IpList.Id,
			"name": resp.IpList.Name,
		}
	}
	if resp.IpItem != nil {
		resultMap["item"] = maps.Map{
			"id":             resp.IpItem.Id,
			"value":          resp.IpItem.Value,
			"ipFrom":         resp.IpItem.IpFrom,
			"ipTo":           resp.IpItem.IpTo,
			"reason":         resp.IpItem.Reason,
			"expiredAt":      resp.IpItem.ExpiredAt,
			"createdTime":    timeutil.FormatTime("Y-m-d", resp.IpItem.CreatedAt),
			"expiredTime":    timeutil.FormatTime("Y-m-d H:i:s", resp.IpItem.ExpiredAt),
			"type":           resp.IpItem.Type,
			"eventLevelName": firewallconfigs.FindFirewallEventLevelName(resp.IpItem.EventLevel),
			"listType":       resp.IpItem.ListType,
		}
	}

	if resp.RegionCountry != nil {
		resultMap["country"] = maps.Map{
			"id":   resp.RegionCountry.Id,
			"name": resp.RegionCountry.Name,
		}
	}

	if resp.RegionProvince != nil {
		resultMap["province"] = maps.Map{
			"id":   resp.RegionProvince.Id,
			"name": resp.RegionProvince.Name,
		}
	}

	this.Data["result"] = resultMap

	this.Success()
}
