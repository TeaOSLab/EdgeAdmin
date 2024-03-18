package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"strings"
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
	grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if grantResp.NodeGrant == nil {
		this.WriteString("can not find the grant")
		return
	}

	// TODO 处理节点专用的认证

	var grant = grantResp.NodeGrant

	var privateKey = grant.PrivateKey
	const maskLength = 64
	if len(privateKey) > maskLength+32 {
		privateKey = privateKey[:maskLength] + strings.Repeat("*", len(privateKey)-maskLength)
	}

	this.Data["grant"] = maps.Map{
		"id":          grant.Id,
		"name":        grant.Name,
		"method":      grant.Method,
		"methodName":  grantutils.FindGrantMethodName(grant.Method, this.LangCode()),
		"username":    grant.Username,
		"password":    strings.Repeat("*", len(grant.Password)),
		"privateKey":  privateKey,
		"passphrase":  strings.Repeat("*", len(grant.Passphrase)),
		"description": grant.Description,
		"su":          grant.Su,
	}

	// 使用此认证的集群
	clusterMaps := []maps.Map{}
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClustersWithNodeGrantId(this.AdminContext(), &pb.FindAllEnabledNodeClustersWithNodeGrantIdRequest{NodeGrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, cluster := range clustersResp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	// 使用此认证的节点
	nodeMaps := []maps.Map{}
	nodesResp, err := this.RPC().NodeRPC().FindAllEnabledNodesWithNodeGrantId(this.AdminContext(), &pb.FindAllEnabledNodesWithNodeGrantIdRequest{NodeGrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, node := range nodesResp.Nodes {
		if node.NodeCluster == nil {
			continue
		}

		clusterMap := maps.Map{
			"id":   node.NodeCluster.Id,
			"name": node.NodeCluster.Name,
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
