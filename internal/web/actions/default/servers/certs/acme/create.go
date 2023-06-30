package acme

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
	// 证书服务商
	providersResp, err := this.RPC().ACMEProviderRPC().FindAllACMEProviders(this.AdminContext(), &pb.FindAllACMEProvidersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var providerMaps = []maps.Map{}
	for _, provider := range providersResp.AcmeProviders {
		providerMaps = append(providerMaps, maps.Map{
			"name": provider.Name,
			"code": provider.Code,
		})
	}
	this.Data["providers"] = providerMaps

	// 域名解析服务商
	dnsProvidersResp, err := this.RPC().DNSProviderRPC().FindAllEnabledDNSProviders(this.AdminContext(), &pb.FindAllEnabledDNSProvidersRequest{
		AdminId: this.AdminId(),
		UserId:  0,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	dnsProviderMaps := []maps.Map{}
	for _, provider := range dnsProvidersResp.DnsProviders {
		dnsProviderMaps = append(dnsProviderMaps, maps.Map{
			"id":       provider.Id,
			"name":     provider.Name,
			"typeName": provider.TypeName,
		})
	}
	this.Data["dnsProviders"] = dnsProviderMaps

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	PlatformUserId int64
	TaskId         int64
	AuthType       string
	AcmeUserId     int64
	DnsProviderId  int64
	DnsDomain      string
	Domains        []string
	AutoRenew      bool
	AuthURL        string

	Must *actions.Must
}) {
	if params.AuthType != "dns" && params.AuthType != "http" {
		this.Fail("无法识别的认证方式'" + params.AuthType + "'")
	}

	if params.AcmeUserId <= 0 {
		this.Fail("请选择一个申请证书的用户")
	}

	// 校验DNS相关信息
	dnsDomain := strings.ToLower(params.DnsDomain)
	if params.AuthType == "dns" {
		if params.DnsProviderId <= 0 {
			this.Fail("请选择DNS服务商")
		}
		if len(params.DnsDomain) == 0 {
			this.Fail("请输入顶级域名")
		}
		if !domainutils.ValidateDomainFormat(dnsDomain) {
			this.Fail("请输入正确的顶级域名")
		}
	}

	if len(params.Domains) == 0 {
		this.Fail("请输入证书域名列表")
	}
	var realDomains = []string{}
	for _, domain := range params.Domains {
		domain = strings.ToLower(domain)
		if params.AuthType == "dns" { // DNS认证
			if !strings.HasSuffix(domain, "."+dnsDomain) && domain != dnsDomain {
				this.Fail("证书域名中的" + domain + "和顶级域名不一致")
			}
		} else if params.AuthType == "http" { // HTTP认证
			if strings.Contains(domain, "*") {
				this.Fail("在HTTP认证时域名" + domain + "不能包含通配符")
			}
		}
		realDomains = append(realDomains, domain)
	}

	if params.TaskId == 0 {
		createResp, err := this.RPC().ACMETaskRPC().CreateACMETask(this.AdminContext(), &pb.CreateACMETaskRequest{
			UserId:        params.PlatformUserId,
			AuthType:      params.AuthType,
			AcmeUserId:    params.AcmeUserId,
			DnsProviderId: params.DnsProviderId,
			DnsDomain:     dnsDomain,
			Domains:       realDomains,
			AutoRenew:     params.AutoRenew,
			AuthURL:       params.AuthURL,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		params.TaskId = createResp.AcmeTaskId
		defer this.CreateLogInfo(codes.ACMETask_LogCreateACMETask, createResp.AcmeTaskId)
	} else {
		_, err := this.RPC().ACMETaskRPC().UpdateACMETask(this.AdminContext(), &pb.UpdateACMETaskRequest{
			AcmeTaskId:    params.TaskId,
			AcmeUserId:    params.AcmeUserId,
			DnsProviderId: params.DnsProviderId,
			DnsDomain:     dnsDomain,
			Domains:       realDomains,
			AutoRenew:     params.AutoRenew,
			AuthURL:       params.AuthURL,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		defer this.CreateLogInfo(codes.ACMETask_LogUpdateACMETask, params.TaskId)
	}

	this.Data["taskId"] = params.TaskId

	this.Success()
}
