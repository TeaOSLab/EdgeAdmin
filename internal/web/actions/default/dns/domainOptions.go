package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

// DomainOptionsAction 域名列表选项
type DomainOptionsAction struct {
	actionutils.ParentAction
}

func (this *DomainOptionsAction) RunPost(params struct {
	ProviderId int64
}) {
	domainsResp, err := this.RPC().DNSDomainRPC().FindAllEnabledBasicDNSDomainsWithDNSProviderId(this.AdminContext(), &pb.FindAllEnabledBasicDNSDomainsWithDNSProviderIdRequest{
		DnsProviderId: params.ProviderId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domainMaps := []maps.Map{}
	for _, domain := range domainsResp.DnsDomains {
		// 未开启或者已删除的先跳过
		if !domain.IsOn || domain.IsDeleted {
			continue
		}

		domainMaps = append(domainMaps, maps.Map{
			"id":   domain.Id,
			"name": domain.Name,
		})
	}
	this.Data["domains"] = domainMaps

	this.Success()
}
