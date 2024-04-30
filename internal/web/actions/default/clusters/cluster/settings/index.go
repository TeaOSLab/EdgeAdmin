package settings

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("basic")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	// 基本信息
	clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cluster := clusterResp.NodeCluster
	if cluster == nil {
		this.WriteString("not found cluster")
		return
	}

	// 认证
	var grantMap interface{} = nil

	if cluster.NodeGrantId > 0 {
		grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: cluster.NodeGrantId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var grant = grantResp.NodeGrant
		if grant != nil {
			grantMap = maps.Map{
				"id":         grant.Id,
				"name":       grant.Name,
				"method":     grant.Method,
				"methodName": grantutils.FindGrantMethodName(grant.Method, this.LangCode()),
			}
		}
	}
	this.Data["grant"] = grantMap

	// 时区
	this.Data["timeZoneGroups"] = nodeconfigs.FindAllTimeZoneGroups()
	this.Data["timeZoneLocations"] = nodeconfigs.FindAllTimeZoneLocations()

	if len(cluster.TimeZone) == 0 {
		cluster.TimeZone = nodeconfigs.DefaultTimeZoneLocation
	}
	this.Data["timeZoneLocation"] = nodeconfigs.FindTimeZoneLocation(cluster.TimeZone)

	// 时钟
	var clockConfig = nodeconfigs.DefaultClockConfig()
	if len(cluster.ClockJSON) > 0 {
		err = json.Unmarshal(cluster.ClockJSON, clockConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if clockConfig == nil {
			clockConfig = nodeconfigs.DefaultClockConfig()
		}
	}

	// SSH参数
	var sshParams = nodeconfigs.DefaultSSHParams()
	if len(cluster.SshParamsJSON) > 0 {
		err = json.Unmarshal(cluster.SshParamsJSON, sshParams)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// DNS信息
	var fullDomainName = ""
	if len(cluster.DnsName) > 0 && cluster.DnsDomainId > 0 {
		domainResp, err := this.RPC().DNSDomainRPC().FindBasicDNSDomain(this.AdminContext(), &pb.FindBasicDNSDomainRequest{DnsDomainId: cluster.DnsDomainId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if domainResp.DnsDomain != nil {
			fullDomainName = cluster.DnsName + "." + domainResp.DnsDomain.Name
		}
	}

	this.Data["cluster"] = maps.Map{
		"id":                  cluster.Id,
		"name":                cluster.Name,
		"installDir":          cluster.InstallDir,
		"timeZone":            cluster.TimeZone,
		"nodeMaxThreads":      cluster.NodeMaxThreads,
		"autoOpenPorts":       cluster.AutoOpenPorts,
		"clock":               clockConfig,
		"autoRemoteStart":     cluster.AutoRemoteStart,
		"autoInstallNftables": cluster.AutoInstallNftables,
		"autoSystemTuning":    cluster.AutoSystemTuning,
		"autoTrimDisks":       cluster.AutoTrimDisks,
		"maxConcurrentReads":  cluster.MaxConcurrentReads,
		"maxConcurrentWrites": cluster.MaxConcurrentWrites,
		"sshParams":           sshParams,
		"domainName":          fullDomainName,
	}

	// 默认值
	this.Data["defaultNodeMaxThreads"] = nodeconfigs.DefaultMaxThreads
	this.Data["defaultNodeMaxThreadsMin"] = nodeconfigs.DefaultMaxThreadsMin
	this.Data["defaultNodeMaxThreadsMax"] = nodeconfigs.DefaultMaxThreadsMax

	this.Show()
}

// RunPost 保存设置
func (this *IndexAction) RunPost(params struct {
	ClusterId           int64
	Name                string
	GrantId             int64
	SshParamsPort       int
	InstallDir          string
	TimeZone            string
	NodeMaxThreads      int32
	AutoOpenPorts       bool
	ClockAutoSync       bool
	ClockServer         string
	ClockCheckChrony    bool
	AutoRemoteStart     bool
	AutoInstallNftables bool
	AutoSystemTuning    bool
	AutoTrimDisks       bool
	MaxConcurrentReads  int32
	MaxConcurrentWrites int32

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLogInfo(codes.NodeCluster_LogUpdateClusterBasicSettings, params.ClusterId)

	params.Must.
		Field("name", params.Name).
		Require("请输入集群名称")

	if params.NodeMaxThreads > 0 {
		params.Must.
			Field("nodeMaxThreads", params.NodeMaxThreads).
			Gte(int64(nodeconfigs.DefaultMaxThreadsMin), "单节点最大线程数最小值不能小于"+types.String(nodeconfigs.DefaultMaxThreadsMin)).
			Lte(int64(nodeconfigs.DefaultMaxThreadsMax), "单节点最大线程数最大值不能大于"+types.String(nodeconfigs.DefaultMaxThreadsMax))
	}

	// ssh
	var sshParams = nodeconfigs.DefaultSSHParams()
	sshParams.Port = params.SshParamsPort
	sshParamsJSON, err := json.Marshal(sshParams)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// clock
	var clockConfig = nodeconfigs.DefaultClockConfig()
	clockConfig.AutoSync = params.ClockAutoSync
	clockConfig.Server = params.ClockServer
	clockConfig.CheckChrony = params.ClockCheckChrony
	clockConfigJSON, err := json.Marshal(clockConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	err = clockConfig.Init()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().NodeClusterRPC().UpdateNodeCluster(this.AdminContext(), &pb.UpdateNodeClusterRequest{
		NodeClusterId:       params.ClusterId,
		Name:                params.Name,
		NodeGrantId:         params.GrantId,
		InstallDir:          params.InstallDir,
		TimeZone:            params.TimeZone,
		NodeMaxThreads:      params.NodeMaxThreads,
		AutoOpenPorts:       params.AutoOpenPorts,
		ClockJSON:           clockConfigJSON,
		AutoRemoteStart:     params.AutoRemoteStart,
		AutoInstallNftables: params.AutoInstallNftables,
		AutoSystemTuning:    params.AutoSystemTuning,
		AutoTrimDisks:       params.AutoTrimDisks,
		SshParamsJSON:       sshParamsJSON,
		MaxConcurrentReads:  params.MaxConcurrentReads,
		MaxConcurrentWrites: params.MaxConcurrentWrites,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
