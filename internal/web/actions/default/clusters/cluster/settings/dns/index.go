package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("dns")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	// 是否有域名可选
	hasDomainsResp, err := this.RPC().DNSDomainRPC().ExistAvailableDomains(this.AdminContext(), &pb.ExistAvailableDomainsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasDomains"] = hasDomainsResp.Exist

	// 当前集群的DNS信息
	this.Data["domainId"] = 0
	this.Data["domainName"] = ""
	this.Data["dnsName"] = ""

	dnsInfoResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterDNS(this.AdminContext(), &pb.FindEnabledNodeClusterDNSRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["dnsName"] = dnsInfoResp.Name
	this.Data["nodesAutoSync"] = dnsInfoResp.NodesAutoSync
	this.Data["serversAutoSync"] = dnsInfoResp.ServersAutoSync

	var domainProviderMap = maps.Map{
		"id":   0,
		"name": "",
	}
	if dnsInfoResp.Domain != nil {
		this.Data["domainId"] = dnsInfoResp.Domain.Id
		this.Data["domainName"] = dnsInfoResp.Domain.Name

		if dnsInfoResp.Provider != nil {
			domainProviderMap = maps.Map{
				"id":   dnsInfoResp.Provider.Id,
				"name": dnsInfoResp.Provider.Name,
			}
		}
	}
	this.Data["domainProvider"] = domainProviderMap

	if len(dnsInfoResp.CnameRecords) == 0 {
		this.Data["cnameRecords"] = []string{}
	} else {
		this.Data["cnameRecords"] = dnsInfoResp.CnameRecords
	}
	this.Data["ttl"] = dnsInfoResp.Ttl
	this.Data["cnameAsDomain"] = dnsInfoResp.CnameAsDomain
	this.Data["includingLnNodes"] = dnsInfoResp.IncludingLnNodes

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId int64

	DnsDomainId      int64
	DnsName          string
	NodesAutoSync    bool
	ServersAutoSync  bool
	CnameRecords     []string
	Ttl              int32
	CnameAsDomain    bool
	IncludingLnNodes bool

	ConfirmResetDomain bool // 是否确认重置域名

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 创建日志
	defer this.CreateLogInfo(codes.DNS_LogUpdateClusterDNS, params.ClusterId)

	if !params.ConfirmResetDomain {
		if params.DnsDomainId <= 0 {
			this.Fail("请选择集群的主域名")
		}

		params.Must.
			Field("dnsName", params.DnsName).
			Require("请输入DNS子域名")
	}

	// 检查DNS名称
	if len(params.DnsName) > 0 {
		if !domainutils.ValidateDomainFormat(params.DnsName) {
			this.FailField("dnsName", "请输入正确的DNS子域名")
		}

		// 检查是否已经被使用
		resp, err := this.RPC().NodeClusterRPC().CheckNodeClusterDNSName(this.AdminContext(), &pb.CheckNodeClusterDNSNameRequest{
			NodeClusterId: params.ClusterId,
			DnsName:       params.DnsName,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if resp.IsUsed {
			this.FailField("dnsName", "此DNS子域名已经被使用，请换一个再试")
		}
	}

	_, err := this.RPC().NodeClusterRPC().UpdateNodeClusterDNS(this.AdminContext(), &pb.UpdateNodeClusterDNSRequest{
		NodeClusterId:    params.ClusterId,
		DnsName:          params.DnsName,
		DnsDomainId:      params.DnsDomainId,
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
