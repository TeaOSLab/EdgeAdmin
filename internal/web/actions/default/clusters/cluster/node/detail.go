package node

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/ipAddresses/ipaddressutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/iputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type DetailAction struct {
	actionutils.ParentAction
}

func (this *DetailAction) Init() {
	this.Nav("", "node", "node")
	this.SecondMenu("nodes")
}

func (this *DetailAction) RunGet(params struct {
	NodeId    int64
	ClusterId int64
}) {
	this.Data["nodeId"] = params.NodeId

	nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var node = nodeResp.Node
	if node == nil {
		this.WriteString("找不到要操作的节点")
		return
	}

	// 主集群
	var clusterMap maps.Map = nil
	if node.NodeCluster != nil {
		var clusterId = node.NodeCluster.Id
		clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{NodeClusterId: clusterId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var cluster = clusterResp.NodeCluster
		if cluster != nil {
			clusterMap = maps.Map{
				"id":         cluster.Id,
				"name":       cluster.Name,
				"installDir": cluster.InstallDir,
			}
		}
	}

	// 从集群
	var secondaryClustersMaps = []maps.Map{}
	for _, cluster := range node.SecondaryNodeClusters {
		secondaryClustersMaps = append(secondaryClustersMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
			"isOn": cluster.IsOn,
		})
	}

	// 当前访问集群的DNS设置
	clusterDNSInfo, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterDNS(this.AdminContext(), &pb.FindEnabledNodeClusterDNSRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["dnsIsExcludingLnNode"] = clusterDNSInfo != nil && !clusterDNSInfo.IncludingLnNodes && node.Level > 1

	// IP地址
	ipAddressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledNodeIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledNodeIPAddressesWithNodeIdRequest{
		NodeId: params.NodeId,
		Role:   nodeconfigs.NodeRoleNode,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var ipAddresses = ipAddressesResp.NodeIPAddresses
	var ipAddressMaps = []maps.Map{}
	for _, addr := range ipAddressesResp.NodeIPAddresses {
		thresholds, err := ipaddressutils.InitNodeIPAddressThresholds(this.Parent(), addr.Id)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 是否有备用IP
		var originIP = addr.Ip
		if len(addr.BackupIP) > 0 {
			addr.Ip = addr.BackupIP
		}

		// 专属集群
		var addrClusterMaps = []maps.Map{}
		for _, addrCluster := range addr.NodeClusters {
			addrClusterMaps = append(addrClusterMaps, maps.Map{
				"id":   addrCluster.Id,
				"name": addrCluster.Name,
			})
		}

		ipAddressMaps = append(ipAddressMaps, maps.Map{
			"id":         addr.Id,
			"name":       addr.Name,
			"originIP":   originIP,
			"ip":         addr.Ip,
			"canAccess":  addr.CanAccess,
			"isOn":       addr.IsOn,
			"isUp":       addr.IsUp,
			"clusters":   addrClusterMaps,
			"thresholds": thresholds,
		})
	}

	// DNS相关
	var clusters = []*pb.NodeCluster{node.NodeCluster}
	clusters = append(clusters, node.SecondaryNodeClusters...)
	var recordMaps = []maps.Map{}
	var routeMaps = []maps.Map{}
	for _, cluster := range clusters {
		dnsInfoResp, err := this.RPC().NodeRPC().FindEnabledNodeDNS(this.AdminContext(), &pb.FindEnabledNodeDNSRequest{
			NodeId:        params.NodeId,
			NodeClusterId: cluster.Id,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var dnsInfo = dnsInfoResp.Node
		if len(dnsInfo.DnsDomainName) == 0 || len(dnsInfo.NodeClusterDNSName) == 0 {
			continue
		}
		var domainName = dnsInfo.DnsDomainName

		// 默认线路
		if len(dnsInfo.Routes) == 0 {
			dnsInfo.Routes = append(dnsInfo.Routes, &pb.DNSRoute{})
		} else {
			for _, route := range dnsInfo.Routes {
				routeMaps = append(routeMaps, maps.Map{
					"domainName": domainName,
					"code":       route.Code,
					"name":       route.Name,
				})
			}
		}

		for _, addr := range ipAddresses {
			if !addr.CanAccess || !addr.IsUp || !addr.IsOn {
				continue
			}

			// 过滤集群
			if len(addr.NodeClusters) > 0 {
				var inCluster = false
				for _, addrCluster := range addr.NodeClusters {
					if addrCluster.Id == cluster.Id {
						inCluster = true
					}
				}
				if !inCluster {
					continue
				}
			}

			for _, route := range dnsInfo.Routes {
				var recordType = "A"
				if iputils.IsIPv6(addr.Ip) {
					recordType = "AAAA"
				}
				recordMaps = append(recordMaps, maps.Map{
					"name":        dnsInfo.NodeClusterDNSName + "." + domainName,
					"type":        recordType,
					"route":       route.Name,
					"value":       addr.Ip,
					"clusterName": cluster.Name,
					"isBackup":    dnsInfo.IsBackupForCluster || dnsInfo.IsBackupForGroup,
					"isOffline":   dnsInfo.IsOffline,
				})
			}
		}
	}

	// 登录信息
	var loginMap maps.Map = nil
	if node.NodeLogin != nil {
		loginParams := maps.Map{}
		if len(node.NodeLogin.Params) > 0 {
			err = json.Unmarshal(node.NodeLogin.Params, &loginParams)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		var grantMap = maps.Map{}
		var grantId = loginParams.GetInt64("grantId")
		if grantId > 0 {
			grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: grantId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if grantResp.NodeGrant != nil {
				grantMap = maps.Map{
					"id":         grantResp.NodeGrant.Id,
					"name":       grantResp.NodeGrant.Name,
					"method":     grantResp.NodeGrant.Method,
					"methodName": grantutils.FindGrantMethodName(grantResp.NodeGrant.Method, this.LangCode()),
					"username":   grantResp.NodeGrant.Username,
				}
			}
		}

		loginMap = maps.Map{
			"id":     node.NodeLogin.Id,
			"name":   node.NodeLogin.Name,
			"type":   node.NodeLogin.Type,
			"params": loginParams,
			"grant":  grantMap,
		}
	}

	// 运行状态
	var status = &nodeconfigs.NodeStatus{}
	this.Data["nodeDatetime"] = ""
	this.Data["nodeTimeDiff"] = 0
	if len(node.StatusJSON) > 0 {
		err = json.Unmarshal(node.StatusJSON, &status)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		status.IsActive = status.IsActive && time.Now().Unix()-status.UpdatedAt <= 60 // N秒之内认为活跃

		if status.Timestamp > 0 {
			this.Data["nodeDatetime"] = timeutil.FormatTime("Y-m-d H:i:s", status.Timestamp)
			if status.UpdatedAt > 0 {
				var diff = status.UpdatedAt - status.Timestamp
				if diff < 0 {
					diff = -diff
				}
				this.Data["nodeTimeDiff"] = diff
			}
		} else if status.UpdatedAt > 0 {
			this.Data["nodeDatetime"] = timeutil.FormatTime("Y-m-d H:i:s", status.UpdatedAt)
		}
	}

	// 检查是否有新版本
	if len(status.OS) > 0 {
		checkVersionResp, err := this.RPC().NodeRPC().CheckNodeLatestVersion(this.AdminContext(), &pb.CheckNodeLatestVersionRequest{
			Os:             status.OS,
			Arch:           status.Arch,
			CurrentVersion: status.BuildVersion,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.Data["shouldUpgrade"] = checkVersionResp.HasNewVersion
		this.Data["newVersion"] = checkVersionResp.NewVersion
	} else {
		this.Data["shouldUpgrade"] = false
		this.Data["newVersion"] = ""
	}

	// 分组
	var groupMap maps.Map = nil
	if node.NodeGroup != nil {
		groupMap = maps.Map{
			"id":   node.NodeGroup.Id,
			"name": node.NodeGroup.Name,
		}
	}

	// 区域
	var regionMap maps.Map = nil
	if node.NodeRegion != nil {
		regionMap = maps.Map{
			"id":   node.NodeRegion.Id,
			"name": node.NodeRegion.Name,
		}
	}

	// 缓存硬盘 & 内存容量
	var maxCacheDiskCapacity maps.Map
	if node.MaxCacheDiskCapacity != nil {
		maxCacheDiskCapacity = maps.Map{
			"count": node.MaxCacheDiskCapacity.Count,
			"unit":  node.MaxCacheDiskCapacity.Unit,
		}
	} else {
		maxCacheDiskCapacity = maps.Map{
			"count": 0,
			"unit":  "gb",
		}
	}

	var maxCacheMemoryCapacity maps.Map
	if node.MaxCacheMemoryCapacity != nil {
		maxCacheMemoryCapacity = maps.Map{
			"count": node.MaxCacheMemoryCapacity.Count,
			"unit":  node.MaxCacheMemoryCapacity.Unit,
		}
	} else {
		maxCacheMemoryCapacity = maps.Map{
			"count": 0,
			"unit":  "gb",
		}
	}

	var lnAddrs = node.LnAddrs
	if lnAddrs == nil {
		lnAddrs = []string{}
	}

	// API节点地址
	var apiNodeAddrStrings = []string{}
	var apiNodeAddrs = []*serverconfigs.NetworkAddressConfig{}
	if len(node.ApiNodeAddrsJSON) > 0 {
		err = json.Unmarshal(node.ApiNodeAddrsJSON, &apiNodeAddrs)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, addr := range apiNodeAddrs {
			if addr.Init() == nil {
				apiNodeAddrStrings = append(apiNodeAddrStrings, addr.FullAddresses()...)
			}
		}
	}

	this.Data["node"] = maps.Map{
		"id":                 node.Id,
		"name":               node.Name,
		"ipAddresses":        ipAddressMaps,
		"cluster":            clusterMap,
		"secondaryClusters":  secondaryClustersMaps,
		"login":              loginMap,
		"installDir":         node.InstallDir,
		"isInstalled":        node.IsInstalled,
		"uniqueId":           node.UniqueId,
		"secret":             node.Secret,
		"maxCPU":             node.MaxCPU,
		"isOn":               node.IsOn,
		"records":            recordMaps,
		"routes":             routeMaps,
		"level":              node.Level,
		"levelInfo":          nodeconfigs.FindNodeLevel(int(node.Level)),
		"lnAddrs":            lnAddrs,
		"enableIPLists":      node.EnableIPLists,
		"apiNodeAddrs":       apiNodeAddrStrings,
		"offlineDay":         node.OfflineDay,
		"isOffline":          len(node.OfflineDay) > 0 && node.OfflineDay < timeutil.Format("Ymd"),
		"isBackupForCluster": node.IsBackupForCluster,
		"isBackupForGroup":   node.IsBackupForGroup,

		"status": maps.Map{
			"isActive":             status.IsActive,
			"updatedAt":            status.UpdatedAt,
			"hostname":             status.Hostname,
			"cpuUsage":             status.CPUUsage,
			"cpuUsageText":         fmt.Sprintf("%.2f%%", status.CPUUsage*100),
			"memUsage":             status.MemoryUsage,
			"memUsageText":         fmt.Sprintf("%.2f%%", status.MemoryUsage*100),
			"connectionCount":      status.ConnectionCount,
			"buildVersion":         status.BuildVersion,
			"cpuPhysicalCount":     status.CPUPhysicalCount,
			"cpuLogicalCount":      status.CPULogicalCount,
			"load1m":               numberutils.FormatFloat2(status.Load1m),
			"load5m":               numberutils.FormatFloat2(status.Load5m),
			"load15m":              numberutils.FormatFloat2(status.Load15m),
			"cacheTotalDiskSize":   numberutils.FormatBytes(status.CacheTotalDiskSize),
			"cacheTotalMemorySize": numberutils.FormatBytes(status.CacheTotalMemorySize),
			"exePath":              status.ExePath,
			"apiSuccessPercent":    status.APISuccessPercent,
			"apiAvgCostSeconds":    status.APIAvgCostSeconds,
			"diskWritingSpeedMB":   status.DiskWritingSpeedMB,
		},

		"group":  groupMap,
		"region": regionMap,

		"maxCacheDiskCapacity":   maxCacheDiskCapacity,
		"maxCacheMemoryCapacity": maxCacheMemoryCapacity,
	}

	this.Show()
}
