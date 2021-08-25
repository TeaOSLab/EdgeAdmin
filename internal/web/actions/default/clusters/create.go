package clusters

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "cluster", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
	// 是否有域名
	hasDomainsResp, err := this.RPC().DNSDomainRPC().ExistAvailableDomains(this.AdminContext(), &pb.ExistAvailableDomainsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasDomains"] = hasDomainsResp.Exist

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name string

	// 缓存策略
	CachePolicyId int64

	// WAF策略
	HttpFirewallPolicyId int64

	// SSH相关
	GrantId            int64
	InstallDir         string
	SystemdServiceIsOn bool

	// DNS相关
	DnsDomainId int64
	DnsName     string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入集群名称")

	// 检查DNS名称
	if len(params.DnsName) > 0 {
		if !domainutils.ValidateDomainFormat(params.DnsName) {
			this.FailField("dnsName", "请输入正确的DNS子域名")
		}

		// 检查是否已经被使用
		resp, err := this.RPC().NodeClusterRPC().CheckNodeClusterDNSName(this.AdminContext(), &pb.CheckNodeClusterDNSNameRequest{
			NodeClusterId: 0,
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

	// TODO 检查DnsDomainId的有效性

	// 系统服务
	systemServices := map[string]interface{}{}
	if params.SystemdServiceIsOn {
		systemServices[nodeconfigs.SystemServiceTypeSystemd] = &nodeconfigs.SystemdServiceConfig{
			IsOn: true,
		}
	}
	systemServicesJSON, err := json.Marshal(systemServices)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	createResp, err := this.RPC().NodeClusterRPC().CreateNodeCluster(this.AdminContext(), &pb.CreateNodeClusterRequest{
		Name:                 params.Name,
		NodeGrantId:          params.GrantId,
		InstallDir:           params.InstallDir,
		DnsDomainId:          params.DnsDomainId,
		DnsName:              params.DnsName,
		HttpCachePolicyId:    params.CachePolicyId,
		HttpFirewallPolicyId: params.HttpFirewallPolicyId,
		SystemServicesJSON:   systemServicesJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "创建节点集群：%d", createResp.NodeClusterId)

	this.Data["clusterId"] = createResp.NodeClusterId

	this.Success()
}
