package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("waf")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	webConfigResp, err := this.RPC().ServerRPC().FindAndInitServerWebConfig(this.AdminContext(), &pb.FindAndInitServerWebRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	webConfig := &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(webConfigResp.Config, webConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["webId"] = webConfig.Id
	this.Data["firewallConfig"] = webConfig.FirewallRef

	// 当前已有策略
	policiesResp, err := this.RPC().HTTPFirewallPolicyRPC().FindAllEnabledHTTPFirewallPolicies(this.AdminContext(), &pb.FindAllEnabledHTTPFirewallPoliciesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyMaps := []maps.Map{}
	for _, p := range policiesResp.FirewallPolicies {
		policyMaps = append(policyMaps, maps.Map{
			"id":   p.Id,
			"name": p.Name,
			"isOn": p.IsOn,
		})
	}
	this.Data["firewallPolicies"] = policyMaps

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId        int64
	FirewallJSON []byte

	Must *actions.Must
}) {
	// TODO 检查配置

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPFirewall(this.AdminContext(), &pb.UpdateHTTPFirewallRequest{
		WebId:        params.WebId,
		FirewallJSON: params.FirewallJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
