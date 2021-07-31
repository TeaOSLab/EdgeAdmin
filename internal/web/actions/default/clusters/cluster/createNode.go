package cluster

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strconv"
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
	for _, address := range ipAddresses {
		addressId := address.GetInt64("id")
		if addressId > 0 {
			_, err = this.RPC().NodeIPAddressRPC().UpdateNodeIPAddressNodeId(this.AdminContext(), &pb.UpdateNodeIPAddressNodeIdRequest{
				AddressId: addressId,
				NodeId:    nodeId,
			})
		} else {
			_, err = this.RPC().NodeIPAddressRPC().CreateNodeIPAddress(this.AdminContext(), &pb.CreateNodeIPAddressRequest{
				NodeId:    nodeId,
				Role:      nodeconfigs.NodeRoleNode,
				Name:      address.GetString("name"),
				Ip:        address.GetString("ip"),
				CanAccess: address.GetBool("canAccess"),
			})
		}
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "创建节点 %d", nodeId)

	this.Success()
}
