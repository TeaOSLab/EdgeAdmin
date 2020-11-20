package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 安装节点
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
	if node.Cluster != nil {
		clusterId := node.Cluster.Id
		clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{ClusterId: clusterId})
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

	this.Data["node"] = maps.Map{
		"id":          node.Id,
		"name":        node.Name,
		"installDir":  node.InstallDir,
		"isInstalled": node.IsInstalled,
		"uniqueId":    node.UniqueId,
		"secret":      node.Secret,
		"cluster":     clusterMap,
	}

	this.Show()
}

// 开始安装
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
