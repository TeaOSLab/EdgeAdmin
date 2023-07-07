package ipadmin

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/regionconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

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

	var deniedProvinceIds = []int64{}
	var allowedProvinceIds = []int64{}
	if policyConfig.Inbound != nil && policyConfig.Inbound.Region != nil {
		deniedProvinceIds = policyConfig.Inbound.Region.DenyProvinceIds
		allowedProvinceIds = policyConfig.Inbound.Region.AllowProvinceIds
	}

	provincesResp, err := this.RPC().RegionProvinceRPC().FindAllRegionProvincesWithRegionCountryId(this.AdminContext(), &pb.FindAllRegionProvincesWithRegionCountryIdRequest{
		RegionCountryId: regionconfigs.RegionChinaId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var deniedProvinceMaps = []maps.Map{}
	var allowedProvinceMaps = []maps.Map{}
	for _, province := range provincesResp.RegionProvinces {
		var provinceMap = maps.Map{
			"id":   province.Id,
			"name": province.DisplayName,
		}
		if lists.ContainsInt64(deniedProvinceIds, province.Id) {
			deniedProvinceMaps = append(deniedProvinceMaps, provinceMap)
		}
		if lists.ContainsInt64(allowedProvinceIds, province.Id) {
			allowedProvinceMaps = append(allowedProvinceMaps, provinceMap)
		}
	}
	this.Data["deniedProvinces"] = deniedProvinceMaps
	this.Data["allowedProvinces"] = allowedProvinceMaps

	// except & only URL Patterns
	this.Data["exceptURLPatterns"] = []*shared.URLPattern{}
	this.Data["onlyURLPatterns"] = []*shared.URLPattern{}
	if policyConfig.Inbound != nil && policyConfig.Inbound.Region != nil {
		if len(policyConfig.Inbound.Region.ProvinceExceptURLPatterns) > 0 {
			this.Data["exceptURLPatterns"] = policyConfig.Inbound.Region.ProvinceExceptURLPatterns
		}
		if len(policyConfig.Inbound.Region.ProvinceOnlyURLPatterns) > 0 {
			this.Data["onlyURLPatterns"] = policyConfig.Inbound.Region.ProvinceOnlyURLPatterns
		}
	}

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
	DenyProvinceIds  []int64
	AllowProvinceIds []int64

	ExceptURLPatternsJSON []byte
	OnlyURLPatternsJSON   []byte

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.WAF_LogUpdateForbiddenProvinces, params.FirewallPolicyId)

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
	policyConfig.Inbound.Region.DenyProvinceIds = params.DenyProvinceIds
	policyConfig.Inbound.Region.AllowProvinceIds = params.AllowProvinceIds

	// 例外URL
	var exceptURLPatterns = []*shared.URLPattern{}
	if len(params.ExceptURLPatternsJSON) > 0 {
		err = json.Unmarshal(params.ExceptURLPatternsJSON, &exceptURLPatterns)
		if err != nil {
			this.Fail("校验例外URL参数失败：" + err.Error())
			return
		}
	}
	policyConfig.Inbound.Region.ProvinceExceptURLPatterns = exceptURLPatterns

	// 限制URL
	var onlyURLPatterns = []*shared.URLPattern{}
	if len(params.OnlyURLPatternsJSON) > 0 {
		err = json.Unmarshal(params.OnlyURLPatternsJSON, &onlyURLPatterns)
		if err != nil {
			this.Fail("校验限制URL参数失败：" + err.Error())
			return
		}
	}
	policyConfig.Inbound.Region.ProvinceOnlyURLPatterns = onlyURLPatterns

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
