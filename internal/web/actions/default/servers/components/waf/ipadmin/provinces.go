package ipadmin

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
	this.Nav("", "", "ipadmin")
}

func (this *ProvincesAction) RunGet(params struct {
	FirewallPolicyId int64
}) {
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
	var selectedProvinceIds = []int64{}
	if policyConfig.Inbound != nil && policyConfig.Inbound.Region != nil {
		selectedProvinceIds = policyConfig.Inbound.Region.DenyProvinceIds
	}

	provincesResp, err := this.RPC().RegionProvinceRPC().FindAllRegionProvincesWithRegionCountryId(this.AdminContext(), &pb.FindAllRegionProvincesWithRegionCountryIdRequest{
		RegionCountryId: int64(ChinaCountryId),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var provinceMaps = []maps.Map{}
	for _, province := range provincesResp.RegionProvinces {
		provinceMaps = append(provinceMaps, maps.Map{
			"id":        province.Id,
			"name":      province.DisplayName,
			"isChecked": lists.ContainsInt64(selectedProvinceIds, province.Id),
		})
	}
	this.Data["provinces"] = provinceMaps

	this.Show()
}

func (this *ProvincesAction) RunPost(params struct {
	FirewallPolicyId int64
	ProvinceIds      []int64

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
