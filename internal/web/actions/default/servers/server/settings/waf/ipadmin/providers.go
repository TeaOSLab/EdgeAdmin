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
)

type ProvidersAction struct {
	actionutils.ParentAction
}

func (this *ProvidersAction) Init() {
	this.Nav("", "setting", "provider")
	this.SecondMenu("waf")
}

func (this *ProvidersAction) RunGet(params struct {
	FirewallPolicyId int64
	ServerId         int64
}) {
	this.Data["featureIsOn"] = true
	this.Data["firewallPolicyId"] = params.FirewallPolicyId
	this.Data["subMenuItem"] = "provider"

	// 当前选中的运营商
	policyConfig, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if policyConfig == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
		return
	}

	var deniedProviderIds = []int64{}
	var allowedProviderIds = []int64{}
	var providerHTML string
	if policyConfig.Inbound != nil && policyConfig.Inbound.Region != nil {
		deniedProviderIds = policyConfig.Inbound.Region.DenyProviderIds
		allowedProviderIds = policyConfig.Inbound.Region.AllowProviderIds
		providerHTML = policyConfig.Inbound.Region.ProviderHTML
	}
	this.Data["providerHTML"] = providerHTML

	providerResp, err := this.RPC().RegionProviderRPC().FindAllRegionProviders(this.AdminContext(), &pb.FindAllRegionProvidersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var deniedProviderMaps = []maps.Map{}
	var allowedProviderMaps = []maps.Map{}
	for _, provider := range providerResp.RegionProviders {
		var providerMap = maps.Map{
			"id":   provider.Id,
			"name": provider.Name,
		}
		if lists.ContainsInt64(deniedProviderIds, provider.Id) {
			deniedProviderMaps = append(deniedProviderMaps, providerMap)
		}
		if lists.ContainsInt64(allowedProviderIds, provider.Id) {
			allowedProviderMaps = append(allowedProviderMaps, providerMap)
		}
	}
	this.Data["deniedProviders"] = deniedProviderMaps
	this.Data["allowedProviders"] = allowedProviderMaps

	// except & only URL Patterns
	this.Data["exceptURLPatterns"] = []*shared.URLPattern{}
	this.Data["onlyURLPatterns"] = []*shared.URLPattern{}
	if policyConfig.Inbound != nil && policyConfig.Inbound.Region != nil {
		if len(policyConfig.Inbound.Region.ProviderExceptURLPatterns) > 0 {
			this.Data["exceptURLPatterns"] = policyConfig.Inbound.Region.ProviderExceptURLPatterns
		}
		if len(policyConfig.Inbound.Region.ProviderOnlyURLPatterns) > 0 {
			this.Data["onlyURLPatterns"] = policyConfig.Inbound.Region.ProviderOnlyURLPatterns
		}
	}

	// WAF是否启用
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["wafIsOn"] = webConfig.FirewallRef != nil && webConfig.FirewallRef.IsOn

	// 获取当前服务所在集群的WAF设置
	clusterFirewallPolicy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if clusterFirewallPolicy != nil {
		this.Data["clusterFirewallPolicy"] = maps.Map{
			"id":       clusterFirewallPolicy.Id,
			"name":     clusterFirewallPolicy.Name,
			"isOn":     clusterFirewallPolicy.IsOn,
			"mode":     clusterFirewallPolicy.Mode,
			"modeInfo": firewallconfigs.FindFirewallMode(clusterFirewallPolicy.Mode),
		}
	} else {
		this.Data["clusterFirewallPolicy"] = nil
	}

	this.Show()
}

func (this *ProvidersAction) RunPost(params struct {
	FirewallPolicyId int64

	DenyProviderIds  []int64
	AllowProviderIds []int64

	ExceptURLPatternsJSON []byte
	OnlyURLPatternsJSON   []byte

	ProviderHTML string

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.WAF_LogUpdateForbiddenProviders, params.FirewallPolicyId)

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
	policyConfig.Inbound.Region.DenyProviderIds = params.DenyProviderIds
	policyConfig.Inbound.Region.AllowProviderIds = params.AllowProviderIds

	// 例外URL
	var exceptURLPatterns = []*shared.URLPattern{}
	if len(params.ExceptURLPatternsJSON) > 0 {
		err = json.Unmarshal(params.ExceptURLPatternsJSON, &exceptURLPatterns)
		if err != nil {
			this.Fail("校验例外URL参数失败：" + err.Error())
			return
		}
	}
	policyConfig.Inbound.Region.ProviderExceptURLPatterns = exceptURLPatterns

	// 限制URL
	var onlyURLPatterns = []*shared.URLPattern{}
	if len(params.OnlyURLPatternsJSON) > 0 {
		err = json.Unmarshal(params.OnlyURLPatternsJSON, &onlyURLPatterns)
		if err != nil {
			this.Fail("校验限制URL参数失败：" + err.Error())
			return
		}
	}
	policyConfig.Inbound.Region.ProviderOnlyURLPatterns = onlyURLPatterns

	// 自定义提示
	if len(params.ProviderHTML) > 32<<10 {
		this.Fail("提示内容长度不能超出32K")
		return
	}
	policyConfig.Inbound.Region.ProviderHTML = params.ProviderHTML

	err = policyConfig.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
		return
	}

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
