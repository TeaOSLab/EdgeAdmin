package node

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/ipAddresses/ipaddressutils"
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
	nodeDNS := dnsInfoResp.Node
	dnsRouteMaps := []maps.Map{}
	if nodeDNS != nil {
		for _, dnsInfo := range nodeDNS.Routes {
			dnsRouteMaps = append(dnsRouteMaps, maps.Map{
				"name": dnsInfo.Name,
				"code": dnsInfo.Code,
			})
		}
	}
	this.Data["dnsRoutes"] = dnsRouteMaps
	this.Data["allDNSRoutes"] = []maps.Map{}
	if nodeDNS != nil {
		this.Data["dnsDomainId"] = nodeDNS.DnsDomainId
	} else {
		this.Data["dnsDomainId"] = 0
	}
	if nodeDNS != nil && nodeDNS.DnsDomainId > 0 {
		routesMaps := []maps.Map{}
		routesResp, err := this.RPC().DNSDomainRPC().FindAllDNSDomainRoutes(this.AdminContext(), &pb.FindAllDNSDomainRoutesRequest{DnsDomainId: dnsInfoResp.Node.DnsDomainId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, route := range routesResp.Routes {
			routesMaps = append(routesMaps, maps.Map{
				"name": route.Name,
				"code": route.Code,
			})
		}
		this.Data["allDNSRoutes"] = routesMaps
	}

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

	// 所有集群
	resp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClusters(this.AdminContext(), &pb.FindAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range resp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	LoginId                    int64
	NodeId                     int64
	GroupId                    int64
	RegionId                   int64
	Name                       string
	IPAddressesJSON            []byte `alias:"ipAddressesJSON"`
	ClusterId                  int64
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
	if params.ClusterId <= 0 {
		this.Fail("请选择所在集群")
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
		NodeId:                 params.NodeId,
		GroupId:                params.GroupId,
		RegionId:               params.RegionId,
		Name:                   params.Name,
		NodeClusterId:          params.ClusterId,
		Login:                  loginInfo,
		MaxCPU:                 params.MaxCPU,
		IsOn:                   params.IsOn,
		DnsDomainId:            params.DnsDomainId,
		DnsRoutes:              dnsRouteCodes,
		MaxCacheDiskCapacity:   pbMaxCacheDiskCapacity,
		MaxCacheMemoryCapacity: pbMaxCacheMemoryCapacity,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 禁用老的IP地址
	_, err = this.RPC().NodeIPAddressRPC().DisableAllIPAddressesWithNodeId(this.AdminContext(), &pb.DisableAllIPAddressesWithNodeIdRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 添加新的IP地址
	err = ipaddressutils.UpdateNodeIPAddresses(this.Parent(), params.NodeId, params.IPAddressesJSON)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
