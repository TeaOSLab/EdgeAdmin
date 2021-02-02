package ipadmin

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

const ChinaCountryId = 1

type ProvincesAction struct {
	actionutils.ParentAction
}

func (this *ProvincesAction) Init() {
	this.Nav("", "setting", "province")
	this.SecondMenu("waf")
}

func (this *ProvincesAction) RunGet(params struct {
	FirewallPolicyId int64
	ServerId         int64
}) {
	this.Data["featureIsOn"] = true
	this.Data["firewallPolicyId"] = params.FirewallPolicyId
	this.Data["subMenuItem"] = "province"

	// 当前选中的省份
	policyConfig, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if policyConfig == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
		return
	}
	selectedProvinceIds := []int64{}
	if policyConfig.Inbound != nil && policyConfig.Inbound.Region != nil {
		selectedProvinceIds = policyConfig.Inbound.Region.DenyProvinceIds
	}

	provincesResp, err := this.RPC().RegionProvinceRPC().FindAllEnabledRegionProvincesWithCountryId(this.AdminContext(), &pb.FindAllEnabledRegionProvincesWithCountryIdRequest{
		CountryId: int64(ChinaCountryId),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	provinceMaps := []maps.Map{}
	for _, province := range provincesResp.Provinces {
		provinceMaps = append(provinceMaps, maps.Map{
			"id":        province.Id,
			"name":      province.Name,
			"isChecked": lists.ContainsInt64(selectedProvinceIds, province.Id),
		})
	}
	this.Data["provinces"] = provinceMaps

	// WAF是否启用
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["wafIsOn"] = webConfig.FirewallRef != nil && webConfig.FirewallRef.IsOn

	this.Show()
}

func (this *ProvincesAction) RunPost(params struct {
	FirewallPolicyId int64
	ProvinceIds      []int64

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "WAF策略 %d 设置禁止访问的省份", params.FirewallPolicyId)

	policyConfig, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if policyConfig == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
		return
	}

	if policyConfig.Inbound == nil {
		policyConfig.Inbound = &firewallconfigs.HTTPFirewallInboundConfig{IsOn: true}
	}
	if policyConfig.Inbound.Region == nil {
		policyConfig.Inbound.Region = &firewallconfigs.HTTPFirewallRegionConfig{
			IsOn: true,
		}
	}
	policyConfig.Inbound.Region.DenyProvinceIds = params.ProvinceIds

	inboundJSON, err := json.Marshal(policyConfig.Inbound)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPFirewallPolicyRPC().UpdateHTTPFirewallInboundConfig(this.AdminContext(), &pb.UpdateHTTPFirewallInboundConfigRequest{
		HttpFirewallPolicyId: params.FirewallPolicyId,
		InboundJSON:          inboundJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
