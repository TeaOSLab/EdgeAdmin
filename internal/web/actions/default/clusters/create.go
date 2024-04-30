package clusters

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
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

	// 菜单：集群总数
	totalResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClusters(this.AdminContext(), &pb.CountAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["totalNodeClusters"] = totalResp.Count

	// 菜单：节点总数
	totalNodesResp, err := this.RPC().NodeRPC().CountAllEnabledNodes(this.AdminContext(), &pb.CountAllEnabledNodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["totalNodes"] = totalNodesResp.Count

	// 随机子域名
	var defaultDNSName = "g" + rands.HexString(6) + ".cdn"
	{
		var b = make([]byte, 3)
		_, err = rand.Read(b)
		if err == nil {
			defaultDNSName = fmt.Sprintf("g%x.cdn", b)
		}
	}
	this.Data["defaultDNSName"] = defaultDNSName
	this.Data["dnsName"] = defaultDNSName

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name string

	// 缓存策略
	CachePolicyId int64

	// WAF策略
	HttpFirewallPolicyId int64

	// 服务配置
	MatchDomainStrictly bool

	// SSH相关
	GrantId             int64
	InstallDir          string
	SystemdServiceIsOn  bool
	AutoInstallNftables bool
	AutoSystemTuning    bool
	AutoTrimDisks       bool
	MaxConcurrentReads  int32
	MaxConcurrentWrites int32

	// DNS相关
	DnsDomainId int64
	DnsName     string
	DnsTTL      int32

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

	// 全局服务配置
	var globalServerConfig = serverconfigs.NewGlobalServerConfig()
	globalServerConfig.HTTPAll.MatchDomainStrictly = params.MatchDomainStrictly
	globalServerConfigJSON, err := json.Marshal(globalServerConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 系统服务
	var systemServices = map[string]any{}
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
		Name:                   params.Name,
		NodeGrantId:            params.GrantId,
		InstallDir:             params.InstallDir,
		DnsDomainId:            params.DnsDomainId,
		DnsName:                params.DnsName,
		DnsTTL:                 params.DnsTTL,
		HttpCachePolicyId:      params.CachePolicyId,
		HttpFirewallPolicyId:   params.HttpFirewallPolicyId,
		SystemServicesJSON:     systemServicesJSON,
		GlobalServerConfigJSON: globalServerConfigJSON,
		AutoInstallNftables:    params.AutoInstallNftables,
		AutoSystemTuning:       params.AutoSystemTuning,
		AutoTrimDisks:          params.AutoTrimDisks,
		MaxConcurrentReads:     params.MaxConcurrentReads,
		MaxConcurrentWrites:    params.MaxConcurrentWrites,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLogInfo(codes.NodeCluster_LogCreateCluster, createResp.NodeClusterId)

	this.Data["clusterId"] = createResp.NodeClusterId

	this.Success()
}
