package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// InstallAction 安装节点
type InstallAction struct {
	actionutils.ParentAction
}

func (this *InstallAction) Init() {
	this.Nav("", "node", "install")
	this.SecondMenu("nodes")
}

func (this *InstallAction) RunGet(params struct {
	NodeId int64
}) {
	this.Data["nodeId"] = params.NodeId

	// 节点
	node, err := nodeutils.InitNodeInfo(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 安装信息
	if node.InstallStatus != nil {
		this.Data["installStatus"] = maps.Map{
			"isRunning":  node.InstallStatus.IsRunning,
			"isFinished": node.InstallStatus.IsFinished,
			"isOk":       node.InstallStatus.IsOk,
			"updatedAt":  node.InstallStatus.UpdatedAt,
			"error":      node.InstallStatus.Error,
		}
	} else {
		this.Data["installStatus"] = nil
	}

	// 集群
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

	// API节点列表
	apiNodesResp, err := this.RPC().APINodeRPC().FindAllEnabledAPINodes(this.AdminContext(), &pb.FindAllEnabledAPINodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	apiNodes := apiNodesResp.Nodes
	apiEndpoints := []string{}
	for _, apiNode := range apiNodes {
		if !apiNode.IsOn {
			continue
		}
		apiEndpoints = append(apiEndpoints, apiNode.AccessAddrs...)
	}
	this.Data["apiEndpoints"] = "\"" + strings.Join(apiEndpoints, "\", \"") + "\""

	var nodeMap = this.Data["node"].(maps.Map)
	nodeMap["installDir"] = node.InstallDir
	nodeMap["isInstalled"] = node.IsInstalled
	nodeMap["uniqueId"] = node.UniqueId
	nodeMap["secret"] = node.Secret
	nodeMap["cluster"] = clusterMap

	this.Show()
}

// RunPost 开始安装
func (this *InstallAction) RunPost(params struct {
	NodeId int64

	Must *actions.Must
}) {
	_, err := this.RPC().NodeRPC().InstallNode(this.AdminContext(), &pb.InstallNodeRequest{
		NodeId: params.NodeId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "安装节点 %d", params.NodeId)

	this.Success()
}
