package node

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
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
	NodeId int64
}) {
	this.Data["nodeId"] = params.NodeId

	nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	node := nodeResp.Node
	if node == nil {
		this.WriteString("找不到要操作的节点")
		return
	}

	var clusterMap maps.Map = nil
	if node.NodeCluster != nil {
		clusterId := node.NodeCluster.Id
		clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{NodeClusterId: clusterId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		cluster := clusterResp.NodeCluster
		if cluster != nil {
			clusterMap = maps.Map{
				"id":         cluster.Id,
				"name":       cluster.Name,
				"installDir": cluster.InstallDir,
			}
		}
	}

	// IP地址
	ipAddressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledIPAddressesWithNodeIdRequest{
		NodeId: params.NodeId,
		Role:   nodeconfigs.NodeRoleNode,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	ipAddressMaps := []maps.Map{}
	for _, addr := range ipAddressesResp.Addresses {
		ipAddressMaps = append(ipAddressMaps, maps.Map{
			"id":        addr.Id,
			"name":      addr.Name,
			"ip":        addr.Ip,
			"canAccess": addr.CanAccess,
		})
	}

	// DNS相关
	dnsInfoResp, err := this.RPC().NodeRPC().FindEnabledNodeDNS(this.AdminContext(), &pb.FindEnabledNodeDNSRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	dnsRouteMaps := []maps.Map{}
	recordName := ""
	recordValue := ""
	if dnsInfoResp.Node != nil {
		recordName = dnsInfoResp.Node.NodeClusterDNSName + "." + dnsInfoResp.Node.DnsDomainName
		recordValue = dnsInfoResp.Node.IpAddr
		for _, dnsInfo := range dnsInfoResp.Node.Routes {
			dnsRouteMaps = append(dnsRouteMaps, maps.Map{
				"name": dnsInfo.Name,
				"code": dnsInfo.Code,
			})
		}
	}
	if len(dnsRouteMaps) == 0 {
		dnsRouteMaps = append(dnsRouteMaps, maps.Map{
			"name": "",
			"code": "",
		})
	}
	this.Data["dnsRoutes"] = dnsRouteMaps
	this.Data["dnsRecordName"] = recordName
	this.Data["dnsRecordValue"] = recordValue

	// 登录信息
	var loginMap maps.Map = nil
	if node.Login != nil {
		loginParams := maps.Map{}
		if len(node.Login.Params) > 0 {
			err = json.Unmarshal(node.Login.Params, &loginParams)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		grantMap := maps.Map{}
		grantId := loginParams.GetInt64("grantId")
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
					"methodName": grantutils.FindGrantMethodName(grantResp.NodeGrant.Method),
				}
			}
		}

		loginMap = maps.Map{
			"id":     node.Login.Id,
			"name":   node.Login.Name,
			"type":   node.Login.Type,
			"params": loginParams,
			"grant":  grantMap,
		}
	}

	// 运行状态
	status := &nodeconfigs.NodeStatus{}
	if len(node.StatusJSON) > 0 {
		err = json.Unmarshal(node.StatusJSON, &status)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		status.IsActive = status.IsActive && time.Now().Unix()-status.UpdatedAt <= 60 // N秒之内认为活跃
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
	var maxCacheDiskCapacity maps.Map = nil
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

	var maxCacheMemoryCapacity maps.Map = nil
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

	this.Data["node"] = maps.Map{
		"id":          node.Id,
		"name":        node.Name,
		"ipAddresses": ipAddressMaps,
		"cluster":     clusterMap,
		"login":       loginMap,
		"installDir":  node.InstallDir,
		"isInstalled": node.IsInstalled,
		"uniqueId":    node.UniqueId,
		"secret":      node.Secret,
		"maxCPU":      node.MaxCPU,
		"isOn":        node.IsOn,

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
			"load1m":               fmt.Sprintf("%.2f", status.Load1m),
			"load5m":               fmt.Sprintf("%.2f", status.Load5m),
			"load15m":              fmt.Sprintf("%.2f", status.Load15m),
			"cacheTotalDiskSize":   numberutils.FormatBytes(status.CacheTotalDiskSize),
			"cacheTotalMemorySize": numberutils.FormatBytes(status.CacheTotalMemorySize),
		},

		"group":  groupMap,
		"region": regionMap,

		"maxCacheDiskCapacity":   maxCacheDiskCapacity,
		"maxCacheMemoryCapacity": maxCacheMemoryCapacity,
	}

	this.Show()
}
