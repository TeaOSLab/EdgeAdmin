package cluster

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strconv"
	"strings"
)

// CreateNodeAction 创建节点
type CreateNodeAction struct {
	actionutils.ParentAction
}

func (this *CreateNodeAction) Init() {
	this.Nav("", "node", "create")
	this.SecondMenu("nodes")
}

func (this *CreateNodeAction) RunGet(params struct {
	ClusterId int64
}) {
	leftMenuItems := []maps.Map{
		{
			"name":     "单个创建",
			"url":      "/clusters/cluster/createNode?clusterId=" + strconv.FormatInt(params.ClusterId, 10),
			"isActive": true,
		},
		{
			"name":     "批量创建",
			"url":      "/clusters/cluster/createBatch?clusterId=" + strconv.FormatInt(params.ClusterId, 10),
			"isActive": false,
		},
	}
	this.Data["leftMenuItems"] = leftMenuItems

	// DNS线路
	clusterDNSResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterDNS(this.AdminContext(), &pb.FindEnabledNodeClusterDNSRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	dnsRouteMaps := []maps.Map{}
	this.Data["dnsDomainId"] = 0
	if clusterDNSResp.Domain != nil {
		domainId := clusterDNSResp.Domain.Id
		this.Data["dnsDomainId"] = domainId
		if domainId > 0 {
			routesResp, err := this.RPC().DNSDomainRPC().FindAllDNSDomainRoutes(this.AdminContext(), &pb.FindAllDNSDomainRoutesRequest{DnsDomainId: domainId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			for _, route := range routesResp.Routes {
				dnsRouteMaps = append(dnsRouteMaps, maps.Map{
					"domainId":   domainId,
					"domainName": clusterDNSResp.Domain.Name,
					"name":       route.Name,
					"code":       route.Code,
				})
			}
		}
	}
	this.Data["dnsRoutes"] = dnsRouteMaps

	// API节点列表
	apiNodesResp, err := this.RPC().APINodeRPC().FindAllEnabledAPINodes(this.AdminContext(), &pb.FindAllEnabledAPINodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	apiNodes := apiNodesResp.ApiNodes
	apiEndpoints := []string{}
	for _, apiNode := range apiNodes {
		if !apiNode.IsOn {
			continue
		}
		apiEndpoints = append(apiEndpoints, apiNode.AccessAddrs...)
	}
	this.Data["apiEndpoints"] = "\"" + strings.Join(apiEndpoints, "\", \"") + "\""

	this.Show()
}

func (this *CreateNodeAction) RunPost(params struct {
	Name            string
	IpAddressesJSON []byte
	ClusterId       int64
	GroupId         int64
	RegionId        int64
	GrantId         int64
	SshHost         string
	SshPort         int

	DnsDomainId   int64
	DnsRoutesJSON []byte

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称")

	if len(params.IpAddressesJSON) == 0 {
		this.Fail("请至少添加一个IP地址")
	}

	// TODO 检查cluster
	if params.ClusterId <= 0 {
		this.Fail("请选择所在集群")
	}

	// IP地址
	ipAddresses := []maps.Map{}
	if len(params.IpAddressesJSON) > 0 {
		err := json.Unmarshal(params.IpAddressesJSON, &ipAddresses)
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
		Id:   0,
		Name: "SSH",
		Type: "ssh",
		Params: maps.Map{
			"grantId": params.GrantId,
			"host":    params.SshHost,
			"port":    params.SshPort,
		}.AsJSON(),
	}

	// 保存
	createResp, err := this.RPC().NodeRPC().CreateNode(this.AdminContext(), &pb.CreateNodeRequest{
		Name:          params.Name,
		NodeClusterId: params.ClusterId,
		NodeGroupId:   params.GroupId,
		NodeRegionId:  params.RegionId,
		NodeLogin:     loginInfo,
		DnsDomainId:   params.DnsDomainId,
		DnsRoutes:     dnsRouteCodes,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	nodeId := createResp.NodeId

	// IP地址
	for _, addr := range ipAddresses {
		addrId := addr.GetInt64("id")
		if addrId > 0 {
			_, err = this.RPC().NodeIPAddressRPC().UpdateNodeIPAddressNodeId(this.AdminContext(), &pb.UpdateNodeIPAddressNodeIdRequest{
				NodeIPAddressId: addrId,
				NodeId:          nodeId,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		} else {
			createResp, err := this.RPC().NodeIPAddressRPC().CreateNodeIPAddress(this.AdminContext(), &pb.CreateNodeIPAddressRequest{
				NodeId:    nodeId,
				Role:      nodeconfigs.NodeRoleNode,
				Name:      addr.GetString("name"),
				Ip:        addr.GetString("ip"),
				CanAccess: addr.GetBool("canAccess"),
				IsUp:      addr.GetBool("isUp"),
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			addrId = createResp.NodeIPAddressId
		}

		// 阈值
		var thresholds = addr.GetSlice("thresholds")
		if len(thresholds) > 0 {
			thresholdsJSON, err := json.Marshal(thresholds)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			_, err = this.RPC().NodeIPAddressThresholdRPC().UpdateAllNodeIPAddressThresholds(this.AdminContext(), &pb.UpdateAllNodeIPAddressThresholdsRequest{
				NodeIPAddressId:             addrId,
				NodeIPAddressThresholdsJSON: thresholdsJSON,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "创建节点 %d", nodeId)

	// 响应数据
	this.Data["nodeId"] = nodeId
	nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: nodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if nodeResp.Node != nil {
		var addresses = []string{}
		for _, addrMap := range ipAddresses {
			addresses = append(addresses, addrMap.GetString("ip"))
		}

		var grantMap maps.Map = nil
		grantId := params.GrantId
		if grantId > 0 {
			grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: grantId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if grantResp.NodeGrant != nil && grantResp.NodeGrant.Id > 0 {
				grantMap = maps.Map{
					"id":         grantResp.NodeGrant.Id,
					"name":       grantResp.NodeGrant.Name,
					"method":     grantResp.NodeGrant.Method,
					"methodName": grantutils.FindGrantMethodName(grantResp.NodeGrant.Method),
					"username":   grantResp.NodeGrant.Username,
				}
			}
		}

		this.Data["node"] = maps.Map{
			"id":        nodeResp.Node.Id,
			"name":      nodeResp.Node.Name,
			"uniqueId":  nodeResp.Node.UniqueId,
			"secret":    nodeResp.Node.Secret,
			"addresses": addresses,
			"login": maps.Map{
				"id":   0,
				"name": "SSH",
				"type": "ssh",
				"params": maps.Map{
					"grantId": params.GrantId,
					"host":    params.SshHost,
					"port":    params.SshPort,
				},
			},
			"grant": grantMap,
		}
	}

	this.Success()
}
