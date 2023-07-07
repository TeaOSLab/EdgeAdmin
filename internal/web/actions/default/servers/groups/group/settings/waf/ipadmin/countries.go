package ipadmin

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type CountriesAction struct {
	actionutils.ParentAction
}

func (this *CountriesAction) Init() {
	this.Nav("", "setting", "country")
	this.SecondMenu("waf")
}

func (this *CountriesAction) RunGet(params struct {
	FirewallPolicyId int64
	ServerId         int64
}) {
	this.Data["featureIsOn"] = true
	this.Data["firewallPolicyId"] = params.FirewallPolicyId

	this.Data["subMenuItem"] = "region"

	// 当前选中的地区
	policyConfig, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if policyConfig == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
		return
	}

	var deniedCountryIds = []int64{}
	var allowedCountryIds = []int64{}
	if policyConfig.Inbound != nil && policyConfig.Inbound.Region != nil {
		deniedCountryIds = policyConfig.Inbound.Region.DenyCountryIds
		allowedCountryIds = policyConfig.Inbound.Region.AllowCountryIds
	}

	countriesResp, err := this.RPC().RegionCountryRPC().FindAllRegionCountries(this.AdminContext(), &pb.FindAllRegionCountriesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var deniesCountryMaps = []maps.Map{}
	var allowedCountryMaps = []maps.Map{}
	for _, country := range countriesResp.RegionCountries {
		var countryMap = maps.Map{
			"id":     country.Id,
			"name":   country.DisplayName,
			"letter": strings.ToUpper(string(country.Pinyin[0][0])),
		}
		if lists.ContainsInt64(deniedCountryIds, country.Id) {
			deniesCountryMaps = append(deniesCountryMaps, countryMap)
		}
		if lists.ContainsInt64(allowedCountryIds, country.Id) {
			allowedCountryMaps = append(allowedCountryMaps, countryMap)
		}
	}
	this.Data["deniedCountries"] = deniesCountryMaps
	this.Data["allowedCountries"] = allowedCountryMaps

	// except & only URL Patterns
	this.Data["exceptURLPatterns"] = []*shared.URLPattern{}
	this.Data["onlyURLPatterns"] = []*shared.URLPattern{}
	if policyConfig.Inbound != nil && policyConfig.Inbound.Region != nil {
		if len(policyConfig.Inbound.Region.CountryExceptURLPatterns) > 0 {
			this.Data["exceptURLPatterns"] = policyConfig.Inbound.Region.CountryExceptURLPatterns
		}
		if len(policyConfig.Inbound.Region.CountryOnlyURLPatterns) > 0 {
			this.Data["onlyURLPatterns"] = policyConfig.Inbound.Region.CountryOnlyURLPatterns
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

func (this *CountriesAction) RunPost(params struct {
	FirewallPolicyId int64
	DenyCountryIds   []int64
	AllowCountryIds  []int64

	ExceptURLPatternsJSON []byte
	OnlyURLPatternsJSON   []byte

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.WAF_LogUpdateForbiddenCountries, params.FirewallPolicyId)

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
	policyConfig.Inbound.Region.DenyCountryIds = params.DenyCountryIds
	policyConfig.Inbound.Region.AllowCountryIds = params.AllowCountryIds

	// 例外URL
	var exceptURLPatterns = []*shared.URLPattern{}
	if len(params.ExceptURLPatternsJSON) > 0 {
		err = json.Unmarshal(params.ExceptURLPatternsJSON, &exceptURLPatterns)
		if err != nil {
			this.Fail("校验例外URL参数失败：" + err.Error())
			return
		}
	}
	policyConfig.Inbound.Region.CountryExceptURLPatterns = exceptURLPatterns

	// 限制URL
	var onlyURLPatterns = []*shared.URLPattern{}
	if len(params.OnlyURLPatternsJSON) > 0 {
		err = json.Unmarshal(params.OnlyURLPatternsJSON, &onlyURLPatterns)
		if err != nil {
			this.Fail("校验限制URL参数失败：" + err.Error())
			return
		}
	}
	policyConfig.Inbound.Region.CountryOnlyURLPatterns = onlyURLPatterns

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
