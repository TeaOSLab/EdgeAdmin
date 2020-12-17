package node

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"time"
)

type NodeAction struct {
	actionutils.ParentAction
}

func (this *NodeAction) Init() {
	this.Nav("", "node", "node")
	this.SecondMenu("nodes")
}

func (this *NodeAction) RunGet(params struct {
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
	if node.Cluster != nil {
		clusterId := node.Cluster.Id
		clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{NodeClusterId: clusterId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		cluster := clusterResp.Cluster
		if cluster != nil {
			clusterMap = maps.Map{
				"id":         cluster.Id,
				"name":       cluster.Name,
				"installDir": cluster.InstallDir,
			}
		}
	}

	// IP地址
	ipAddressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledIPAddressesWithNodeIdRequest{NodeId: params.NodeId})
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
	if dnsInfoResp.Node != nil {
		for _, dnsInfo := range dnsInfoResp.Node.Routes {
			dnsRouteMaps = append(dnsRouteMaps, maps.Map{
				"name": dnsInfo.Name,
				"code": dnsInfo.Code,
			})
		}
	}
	this.Data["dnsRoutes"] = dnsRouteMaps

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
			grantResp, err := this.RPC().NodeGrantRPC().FindEnabledGrant(this.AdminContext(), &pb.FindEnabledGrantRequest{GrantId: grantId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if grantResp.Grant != nil {
				grantMap = maps.Map{
					"id":         grantResp.Grant.Id,
					"name":       grantResp.Grant.Name,
					"method":     grantResp.Grant.Method,
					"methodName": grantutils.FindGrantMethodName(grantResp.Grant.Method),
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

	// 分组
	var groupMap maps.Map = nil
	if node.Group != nil {
		groupMap = maps.Map{
			"id":   node.Group.Id,
			"name": node.Group.Name,
		}
	}

	// 区域
	var regionMap maps.Map = nil
	if node.Region != nil {
		regionMap = maps.Map{
			"id":   node.Region.Id,
			"name": node.Region.Name,
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
			"isActive":        status.IsActive,
			"updatedAt":       status.UpdatedAt,
			"hostname":        status.Hostname,
			"cpuUsage":        status.CPUUsage,
			"cpuUsageText":    fmt.Sprintf("%.2f%%", status.CPUUsage*100),
			"memUsage":        status.MemoryUsage,
			"memUsageText":    fmt.Sprintf("%.2f%%", status.MemoryUsage*100),
			"connectionCount": status.ConnectionCount,
		},

		"group":  groupMap,
		"region": regionMap,
	}

	this.Show()
}
