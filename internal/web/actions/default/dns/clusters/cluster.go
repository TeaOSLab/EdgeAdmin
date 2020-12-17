package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type ClusterAction struct {
	actionutils.ParentAction
}

func (this *ClusterAction) Init() {
	this.Nav("", "", "")
}

func (this *ClusterAction) RunGet(params struct {
	ClusterId int64
}) {
	// 集群信息
	clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cluster := clusterResp.Cluster
	if cluster == nil {
		this.NotFound("nodeCluster", params.ClusterId)
		return
	}
	this.Data["cluster"] = maps.Map{
		"id":   cluster.Id,
		"name": cluster.Name,
	}

	// DNS信息
	dnsResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterDNS(this.AdminContext(), &pb.FindEnabledNodeClusterDNSRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	dnsMap := maps.Map{
		"dnsName":          dnsResp.Name,
		"domainId":         0,
		"domainName":       "",
		"providerId":       0,
		"providerName":     "",
		"providerTypeName": "",
	}
	if dnsResp.Domain != nil {
		dnsMap["domainId"] = dnsResp.Domain.Id
		dnsMap["domainName"] = dnsResp.Domain.Name
	}
	if dnsResp.Provider != nil {
		dnsMap["providerId"] = dnsResp.Provider.Id
		dnsMap["providerName"] = dnsResp.Provider.Name
		dnsMap["providerTypeName"] = dnsResp.Provider.TypeName
	}

	this.Data["dnsInfo"] = dnsMap

	// 节点DNS解析记录
	nodesResp, err := this.RPC().NodeRPC().FindAllEnabledNodesDNSWithClusterId(this.AdminContext(), &pb.FindAllEnabledNodesDNSWithClusterIdRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	nodeMaps := []maps.Map{}
	for _, node := range nodesResp.Nodes {
		if len(node.Routes) > 0 {
			for _, route := range node.Routes {
				nodeMaps = append(nodeMaps, maps.Map{
					"id":     node.Id,
					"name":   node.Name,
					"ipAddr": node.IpAddr,
					"route": maps.Map{
						"name": route.Name,
						"code": route.Code,
					},
					"clusterId": node.NodeClusterId,
				})
			}
		} else {
			nodeMaps = append(nodeMaps, maps.Map{
				"id":     node.Id,
				"name":   node.Name,
				"ipAddr": node.IpAddr,
				"route": maps.Map{
					"name": "",
					"code": "",
				},
				"clusterId": node.NodeClusterId,
			})
		}
	}
	this.Data["nodes"] = nodeMaps

	// 代理服务解析记录
	serversResp, err := this.RPC().ServerRPC().FindAllEnabledServersDNSWithClusterId(this.AdminContext(), &pb.FindAllEnabledServersDNSWithClusterIdRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	serverMaps := []maps.Map{}
	for _, server := range serversResp.Servers {
		serverMaps = append(serverMaps, maps.Map{
			"id":      server.Id,
			"name":    server.Name,
			"dnsName": server.DnsName,
		})
	}
	this.Data["servers"] = serverMaps

	// 检查解析记录是否有变化
	checkChangesResp, err := this.RPC().NodeClusterRPC().CheckNodeClusterDNSChanges(this.AdminContext(), &pb.CheckNodeClusterDNSChangesRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["dnsHasChanges"] = checkChangesResp.IsChanged

	this.Show()
}
