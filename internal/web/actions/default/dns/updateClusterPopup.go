package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

// UpdateClusterPopupAction 修改集群的DNS设置
type UpdateClusterPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateClusterPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateClusterPopupAction) RunGet(params struct {
	ClusterId int64
}) {
	this.Data["clusterId"] = params.ClusterId

	dnsInfoResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterDNS(this.AdminContext(), &pb.FindEnabledNodeClusterDNSRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["dnsName"] = dnsInfoResp.Name
	this.Data["nodesAutoSync"] = dnsInfoResp.NodesAutoSync
	this.Data["serversAutoSync"] = dnsInfoResp.ServersAutoSync
	if dnsInfoResp.Domain != nil {
		this.Data["domainId"] = dnsInfoResp.Domain.Id
		this.Data["domain"] = dnsInfoResp.Domain.Name
	} else {
		this.Data["domainId"] = 0
		this.Data["domain"] = ""
	}
	if dnsInfoResp.Provider != nil {
		this.Data["providerType"] = dnsInfoResp.Provider.Type
		this.Data["providerId"] = dnsInfoResp.Provider.Id
	} else {
		this.Data["providerType"] = ""
		this.Data["providerId"] = 0
	}

	if len(dnsInfoResp.CnameRecords) == 0 {
		this.Data["cnameRecords"] = []string{}
	} else {
		this.Data["cnameRecords"] = dnsInfoResp.CnameRecords
	}

	this.Data["ttl"] = dnsInfoResp.Ttl
	this.Data["cnameAsDomain"] = dnsInfoResp.CnameAsDomain
	this.Data["includingLnNodes"] = dnsInfoResp.IncludingLnNodes

	// 所有服务商
	providerTypesResp, err := this.RPC().DNSProviderRPC().FindAllDNSProviderTypes(this.AdminContext(), &pb.FindAllDNSProviderTypesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	providerTypeMaps := []maps.Map{}
	for _, providerType := range providerTypesResp.ProviderTypes {
		providerTypeMaps = append(providerTypeMaps, maps.Map{
			"name": providerType.Name,
			"code": providerType.Code,
		})
	}
	this.Data["providerTypes"] = providerTypeMaps

	this.Show()
}

func (this *UpdateClusterPopupAction) RunPost(params struct {
	ClusterId        int64
	DnsName          string
	DomainId         int64
	NodesAutoSync    bool
	ServersAutoSync  bool
	CnameRecords     []string
	Ttl              int32
	CnameAsDomain    bool
	IncludingLnNodes bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 日志
	defer this.CreateLogInfo(codes.DNS_LogUpdateClusterDNS, params.ClusterId)

	params.Must.
		Field("dnsName", params.DnsName).
		Require("请输入子域名")

	if !domainutils.ValidateDomainFormat(params.DnsName) {
		this.FailField("dnsName", "子域名格式错误")
	}

	checkResp, err := this.RPC().NodeClusterRPC().CheckNodeClusterDNSName(this.AdminContext(), &pb.CheckNodeClusterDNSNameRequest{
		NodeClusterId: params.ClusterId,
		DnsName:       params.DnsName,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if checkResp.IsUsed {
		this.FailField("dnsName", "此子域名已经被占用，请修改后重新提交")
	}

	_, err = this.RPC().NodeClusterRPC().UpdateNodeClusterDNS(this.AdminContext(), &pb.UpdateNodeClusterDNSRequest{
		NodeClusterId:    params.ClusterId,
		DnsName:          params.DnsName,
		DnsDomainId:      params.DomainId,
		NodesAutoSync:    params.NodesAutoSync,
		ServersAutoSync:  params.ServersAutoSync,
		CnameRecords:     params.CnameRecords,
		Ttl:              params.Ttl,
		CnameAsDomain:    params.CnameAsDomain,
		IncludingLnNodes: params.IncludingLnNodes,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
