package acme

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type UpdateTaskPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateTaskPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateTaskPopupAction) RunGet(params struct {
	TaskId int64
}) {
	taskResp, err := this.RPC().ACMETaskRPC().FindEnabledACMETask(this.AdminContext(), &pb.FindEnabledACMETaskRequest{AcmeTaskId: params.TaskId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var task = taskResp.AcmeTask
	if task == nil {
		this.NotFound("acmeTask", params.TaskId)
		return
	}

	var dnsProviderMap maps.Map
	if task.DnsProvider != nil {
		dnsProviderMap = maps.Map{
			"id": task.DnsProvider.Id,
		}
	} else {
		dnsProviderMap = maps.Map{
			"id": 0,
		}
	}

	var acmeUserMap maps.Map
	if task.AcmeUser != nil {
		acmeUserMap = maps.Map{
			"id": task.AcmeUser.Id,
		}
	} else {
		acmeUserMap = maps.Map{
			"id": 0,
		}
	}

	this.Data["task"] = maps.Map{
		"id":          task.Id,
		"authType":    task.AuthType,
		"acmeUser":    acmeUserMap,
		"dnsDomain":   task.DnsDomain,
		"domains":     task.Domains,
		"autoRenew":   task.AutoRenew,
		"isOn":        task.IsOn,
		"authURL":     task.AuthURL,
		"dnsProvider": dnsProviderMap,
	}

	// 域名解析服务商
	providersResp, err := this.RPC().DNSProviderRPC().FindAllEnabledDNSProviders(this.AdminContext(), &pb.FindAllEnabledDNSProvidersRequest{
		AdminId: this.AdminId(),
		UserId:  0,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var providerMaps = []maps.Map{}
	for _, provider := range providersResp.DnsProviders {
		providerMaps = append(providerMaps, maps.Map{
			"id":       provider.Id,
			"name":     provider.Name,
			"typeName": provider.TypeName,
		})
	}
	this.Data["providers"] = providerMaps

	this.Show()
}

func (this *UpdateTaskPopupAction) RunPost(params struct {
	TaskId        int64
	AuthType      string
	AcmeUserId    int64
	DnsProviderId int64
	DnsDomain     string
	DomainsJSON   []byte
	AutoRenew     bool
	AuthURL       string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.ACMETask_LogUpdateACMETask, params.TaskId)

	if params.AuthType != "dns" && params.AuthType != "http" {
		this.Fail("无法识别的认证方式'" + params.AuthType + "'")
	}

	if params.AcmeUserId <= 0 {
		this.Fail("请选择一个申请证书的用户")
	}

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

	var domains = []string{}
	if len(params.DomainsJSON) > 0 {
		err := json.Unmarshal(params.DomainsJSON, &domains)
		if err != nil {
			this.Fail("解析域名数据失败：" + err.Error())
			return
		}
	}

	if len(domains) == 0 {
		this.Fail("请输入证书域名列表")
	}
	var realDomains = []string{}
	for _, domain := range domains {
		domain = strings.ToLower(domain)
		if params.AuthType == "dns" {
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

	this.Success()
}
