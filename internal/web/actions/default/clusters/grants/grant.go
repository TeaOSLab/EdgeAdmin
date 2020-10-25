package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type GrantAction struct {
	actionutils.ParentAction
}

func (this *GrantAction) Init() {
	this.Nav("", "grant", "index")
}

func (this *GrantAction) RunGet(params struct {
	GrantId int64
}) {
	grantResp, err := this.RPC().NodeGrantRPC().FindEnabledGrant(this.AdminContext(), &pb.FindEnabledGrantRequest{GrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if grantResp.Grant == nil {
		this.WriteString("can not find the grant")
		return
	}

	// TODO 处理节点专用的认证

	grant := grantResp.Grant
	this.Data["grant"] = maps.Map{
		"id":          grant.Id,
		"name":        grant.Name,
		"method":      grant.Method,
		"methodName":  grantutils.FindGrantMethodName(grant.Method),
		"username":    grant.Username,
		"password":    grant.Password,
		"privateKey":  grant.PrivateKey,
		"description": grant.Description,
		"su":          grant.Su,
	}

	// 使用此认证的集群
	clusterMaps := []maps.Map{}
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClustersWithGrantId(this.AdminContext(), &pb.FindAllEnabledNodeClustersWithGrantIdRequest{GrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, cluster := range clustersResp.Clusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	// 使用此认证的节点
	nodeMaps := []maps.Map{}
	nodesResp, err := this.RPC().NodeRPC().FindAllEnabledNodesWithGrantId(this.AdminContext(), &pb.FindAllEnabledNodesWithGrantIdRequest{GrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, node := range nodesResp.Nodes {
		if node.Cluster == nil {
			continue
		}

		clusterMap := maps.Map{
			"id":   node.Cluster.Id,
			"name": node.Cluster.Name,
		}

		nodeMaps = append(nodeMaps, maps.Map{
			"id":      node.Id,
			"name":    node.Name,
			"cluster": clusterMap,
			"isOn":    node.IsOn,
		})
	}
	this.Data["nodes"] = nodeMaps

	this.Show()
}
