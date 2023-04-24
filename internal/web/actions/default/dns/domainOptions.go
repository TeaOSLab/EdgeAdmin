package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"sort"
)

// DomainOptionsAction 域名列表选项
type DomainOptionsAction struct {
	actionutils.ParentAction
}

func (this *DomainOptionsAction) RunPost(params struct {
	ProviderId int64
}) {
	domainsResp, err := this.RPC().DNSDomainRPC().FindAllBasicDNSDomainsWithDNSProviderId(this.AdminContext(), &pb.FindAllBasicDNSDomainsWithDNSProviderIdRequest{
		DnsProviderId: params.ProviderId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 排序
	if len(domainsResp.DnsDomains) > 0 {
		sort.Slice(domainsResp.DnsDomains, func(i, j int) bool {
			return domainsResp.DnsDomains[i].Name < domainsResp.DnsDomains[j].Name
		})
	}

	var domainMaps = []maps.Map{}
	for _, domain := range domainsResp.DnsDomains {
		// 未开启或者已删除的先跳过
		if !domain.IsOn || domain.IsDeleted || !domain.IsUp {
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
