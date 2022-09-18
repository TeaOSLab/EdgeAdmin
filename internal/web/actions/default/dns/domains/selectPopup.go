package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectPopupAction) RunGet(params struct {
	DomainId int64
}) {
	this.Data["domainId"] = 0
	this.Data["domainName"] = ""
	this.Data["providerId"] = 0
	this.Data["providerType"] = ""

	// 域名信息
	if params.DomainId > 0 {
		domainResp, err := this.RPC().DNSDomainRPC().FindDNSDomain(this.AdminContext(), &pb.FindDNSDomainRequest{DnsDomainId: params.DomainId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var domain = domainResp.DnsDomain
		if domain != nil {
			this.Data["domainId"] = domain.Id
			this.Data["domainName"] = domain.Name
			this.Data["providerId"] = domain.ProviderId

			providerResp, err := this.RPC().DNSProviderRPC().FindEnabledDNSProvider(this.AdminContext(), &pb.FindEnabledDNSProviderRequest{DnsProviderId: domain.ProviderId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if providerResp.DnsProvider != nil {
				this.Data["providerType"] = providerResp.DnsProvider.Type
			}
		}
	}

	// 所有服务商
	providerTypesResp, err := this.RPC().DNSProviderRPC().FindAllDNSProviderTypes(this.AdminContext(), &pb.FindAllDNSProviderTypesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var providerTypeMaps = []maps.Map{}
	for _, providerType := range providerTypesResp.ProviderTypes {
		providerTypeMaps = append(providerTypeMaps, maps.Map{
			"name": providerType.Name,
			"code": providerType.Code,
		})
	}
	this.Data["providerTypes"] = providerTypeMaps

	this.Show()
}

func (this *SelectPopupAction) RunPost(params struct {
	DomainId int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	this.Data["domainId"] = params.DomainId
	this.Data["domainName"] = ""
	this.Data["providerName"] = ""

	if params.DomainId > 0 {
		domainResp, err := this.RPC().DNSDomainRPC().FindDNSDomain(this.AdminContext(), &pb.FindDNSDomainRequest{DnsDomainId: params.DomainId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if domainResp.DnsDomain != nil {
			this.Data["domainName"] = domainResp.DnsDomain.Name

			// 服务商名称
			var providerId = domainResp.DnsDomain.ProviderId
			if providerId > 0 {
				providerResp, err := this.RPC().DNSProviderRPC().FindEnabledDNSProvider(this.AdminContext(), &pb.FindEnabledDNSProviderRequest{DnsProviderId: providerId})
				if err != nil {
					this.ErrorPage(err)
					return
				}
				if providerResp.DnsProvider != nil {
					this.Data["providerName"] = providerResp.DnsProvider.Name
				}
			}
		} else {
			this.Data["domainId"] = 0
		}
	}

	this.Success()
}
