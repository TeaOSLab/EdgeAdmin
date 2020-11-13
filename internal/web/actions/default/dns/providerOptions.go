package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

// 服务商选项
type ProviderOptionsAction struct {
	actionutils.ParentAction
}

func (this *ProviderOptionsAction) RunPost(params struct {
	Type string
}) {
	providersResp, err := this.RPC().DNSProviderRPC().FindAllEnabledDNSProvidersWithType(this.AdminContext(), &pb.FindAllEnabledDNSProvidersWithTypeRequest{ProviderTypeCode: params.Type})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	providerMaps := []maps.Map{}
	for _, provider := range providersResp.DnsProviders {
		providerMaps = append(providerMaps, maps.Map{
			"id":   provider.Id,
			"name": provider.Name,
		})
	}
	this.Data["providers"] = providerMaps

	this.Success()
}
