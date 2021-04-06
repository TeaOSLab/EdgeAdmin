package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type InstallAction struct {
	actionutils.ParentAction
}

func (this *InstallAction) Init() {
	this.Nav("", "", "install")
}

func (this *InstallAction) RunGet(params struct {
	NodeId int64
}) {
	// 监控节点信息
	nodeResp, err := this.RPC().MonitorNodeRPC().FindEnabledMonitorNode(this.AdminContext(), &pb.FindEnabledMonitorNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	node := nodeResp.Node
	if node == nil {
		this.NotFound("monitorNode", params.NodeId)
		return
	}

	this.Data["node"] = maps.Map{
		"id":       node.Id,
		"name":     node.Name,
		"uniqueId": node.UniqueId,
		"secret":   node.Secret,
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

	this.Show()
}
