package node

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/ipAddresses/ipaddressutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "node", "update")
	this.SecondMenu("nodes")
}

func (this *UpdateAction) RunGet(params struct {
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
		clusterMap = maps.Map{
			"id":   node.NodeCluster.Id,
			"name": node.NodeCluster.Name,
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
	var clusters = []*pb.NodeCluster{node.NodeCluster}
	clusters = append(clusters, node.SecondaryNodeClusters...)
	var allDNSRouteMaps = map[int64][]maps.Map{} // domain id => routes
	var routeMaps = map[int64][]maps.Map{}       // domain id => routes
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
		if dnsInfo.DnsDomainId <= 0 || len(dnsInfo.DnsDomainName) == 0 {
			continue
		}
		var domainId = dnsInfo.DnsDomainId
		var domainName = dnsInfo.DnsDomainName
		if len(dnsInfo.Routes) > 0 {
			for _, route := range dnsInfo.Routes {
				routeMaps[domainId] = append(routeMaps[domainId], maps.Map{
					"domainId":   domainId,
					"domainName": domainName,
					"code":       route.Code,
					"name":       route.Name,
				})
			}
		}

		// 所有线路选项
		routesResp, err := this.RPC().DNSDomainRPC().FindAllDNSDomainRoutes(this.AdminContext(), &pb.FindAllDNSDomainRoutesRequest{DnsDomainId: dnsInfoResp.Node.DnsDomainId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, route := range routesResp.Routes {
			allDNSRouteMaps[domainId] = append(allDNSRouteMaps[domainId], maps.Map{
				"domainId":   domainId,
				"domainName": domainName,
				"name":       route.Name,
				"code":       route.Code,
			})
		}
	}

	var domainRoutes = []maps.Map{}
	for _, m := range routeMaps {
		domainRoutes = append(domainRoutes, m...)
	}
	this.Data["dnsRoutes"] = domainRoutes

	var allDomainRoutes = []maps.Map{}
	for _, m := range allDNSRouteMaps {
		allDomainRoutes = append(allDomainRoutes, m...)
	}
	this.Data["allDNSRoutes"] = allDomainRoutes

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
			"id":     node.NodeLogin.Id,
			"name":   node.NodeLogin.Name,
			"type":   node.NodeLogin.Type,
			"params": loginParams,
			"grant":  grantMap,
		}
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

	var m = maps.Map{
		"id":                     node.Id,
		"name":                   node.Name,
		"ipAddresses":            ipAddressMaps,
		"cluster":                clusterMap,
		"login":                  loginMap,
		"maxCPU":                 node.MaxCPU,
		"isOn":                   node.IsOn,
		"group":                  groupMap,
		"region":                 regionMap,
		"maxCacheDiskCapacity":   maxCacheDiskCapacity,
		"maxCacheMemoryCapacity": maxCacheMemoryCapacity,
	}

	if node.NodeCluster != nil {
		m["primaryCluster"] = maps.Map{
			"id":   node.NodeCluster.Id,
			"name": node.NodeCluster.Name,
		}
	} else {
		m["primaryCluster"] = nil
	}

	if len(node.SecondaryNodeClusters) > 0 {
		var secondaryClusterMaps = []maps.Map{}
		for _, cluster := range node.SecondaryNodeClusters {
			secondaryClusterMaps = append(secondaryClusterMaps, maps.Map{
				"id":   cluster.Id,
				"name": cluster.Name,
			})
		}
		m["secondaryClusters"] = secondaryClusterMaps
	} else {
		m["secondaryClusters"] = []interface{}{}
	}

	this.Data["node"] = m

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	LoginId                    int64
	NodeId                     int64
	GroupId                    int64
	RegionId                   int64
	Name                       string
	IPAddressesJSON            []byte `alias:"ipAddressesJSON"`
	PrimaryClusterId           int64
	SecondaryClusterIds        []byte
	GrantId                    int64
	SshHost                    string
	SshPort                    int
	MaxCPU                     int32
	IsOn                       bool
	MaxCacheDiskCapacityJSON   []byte
	MaxCacheMemoryCapacityJSON []byte

	DnsDomainId   int64
	DnsRoutesJSON []byte

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改节点 %d", params.NodeId)

	if params.NodeId <= 0 {
		this.Fail("要操作的节点不存在")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称")

	// TODO 检查cluster
	if params.PrimaryClusterId <= 0 {
		this.Fail("请选择节点所在主集群")
	}

	var secondaryClusterIds = []int64{}
	if len(params.SecondaryClusterIds) > 0 {
		err := json.Unmarshal(params.SecondaryClusterIds, &secondaryClusterIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// IP地址
	ipAddresses := []maps.Map{}
	if len(params.IPAddressesJSON) > 0 {
		err := json.Unmarshal(params.IPAddressesJSON, &ipAddresses)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	if len(ipAddresses) == 0 {
		this.Fail("请至少输入一个IP地址")
	}

	dnsRouteCodes := []string{}
	if len(params.DnsRoutesJSON) > 0 {
		err := json.Unmarshal(params.DnsRoutesJSON, &dnsRouteCodes)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// TODO 检查登录授权
	loginInfo := &pb.NodeLogin{
		Id:   params.LoginId,
		Name: "SSH",
		Type: "ssh",
		Params: maps.Map{
			"grantId": params.GrantId,
			"host":    params.SshHost,
			"port":    params.SshPort,
		}.AsJSON(),
	}

	// 缓存硬盘 & 内存容量
	var pbMaxCacheDiskCapacity *pb.SizeCapacity
	if len(params.MaxCacheDiskCapacityJSON) > 0 {
		var sizeCapacity = &shared.SizeCapacity{}
		err := json.Unmarshal(params.MaxCacheDiskCapacityJSON, sizeCapacity)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		pbMaxCacheDiskCapacity = &pb.SizeCapacity{
			Count: sizeCapacity.Count,
			Unit:  sizeCapacity.Unit,
		}
	}

	var pbMaxCacheMemoryCapacity *pb.SizeCapacity
	if len(params.MaxCacheMemoryCapacityJSON) > 0 {
		var sizeCapacity = &shared.SizeCapacity{}
		err := json.Unmarshal(params.MaxCacheMemoryCapacityJSON, sizeCapacity)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		pbMaxCacheMemoryCapacity = &pb.SizeCapacity{
			Count: sizeCapacity.Count,
			Unit:  sizeCapacity.Unit,
		}
	}

	// 保存
	_, err := this.RPC().NodeRPC().UpdateNode(this.AdminContext(), &pb.UpdateNodeRequest{
		NodeId:                  params.NodeId,
		NodeGroupId:             params.GroupId,
		NodeRegionId:            params.RegionId,
		Name:                    params.Name,
		NodeClusterId:           params.PrimaryClusterId,
		SecondaryNodeClusterIds: secondaryClusterIds,
		NodeLogin:               loginInfo,
		MaxCPU:                  params.MaxCPU,
		IsOn:                    params.IsOn,
		DnsDomainId:             params.DnsDomainId,
		DnsRoutes:               dnsRouteCodes,
		MaxCacheDiskCapacity:    pbMaxCacheDiskCapacity,
		MaxCacheMemoryCapacity:  pbMaxCacheMemoryCapacity,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 禁用老的IP地址
	_, err = this.RPC().NodeIPAddressRPC().DisableAllIPAddressesWithNodeId(this.AdminContext(), &pb.DisableAllIPAddressesWithNodeIdRequest{
		NodeId: params.NodeId,
		Role:   nodeconfigs.NodeRoleNode,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 添加新的IP地址
	err = ipaddressutils.UpdateNodeIPAddresses(this.Parent(), params.NodeId, nodeconfigs.NodeRoleNode, params.IPAddressesJSON)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
