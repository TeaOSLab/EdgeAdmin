package cluster

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpgradeRemoteAction struct {
	actionutils.ParentAction
}

func (this *UpgradeRemoteAction) Init() {
	this.Nav("", "node", "install")
	this.SecondMenu("nodes")
}

func (this *UpgradeRemoteAction) RunGet(params struct {
	ClusterId int64
}) {
	this.Data["leftMenuItems"] = LeftMenuItemsForInstall(this.AdminContext(), params.ClusterId, "upgrade", this.LangCode())

	nodes := []maps.Map{}
	resp, err := this.RPC().NodeRPC().FindAllUpgradeNodesWithNodeClusterId(this.AdminContext(), &pb.FindAllUpgradeNodesWithNodeClusterIdRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, node := range resp.Nodes {
		loginParams := maps.Map{}
		if node.Node.NodeLogin != nil && len(node.Node.NodeLogin.Params) > 0 {
			err := json.Unmarshal(node.Node.NodeLogin.Params, &loginParams)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		nodes = append(nodes, maps.Map{
			"id":            node.Node.Id,
			"name":          node.Node.Name,
			"os":            node.Os,
			"arch":          node.Arch,
			"oldVersion":    node.OldVersion,
			"newVersion":    node.NewVersion,
			"login":         node.Node.NodeLogin,
			"loginParams":   loginParams,
			"addresses":     node.Node.IpAddresses,
			"installStatus": node.Node.InstallStatus,
		})
	}
	this.Data["nodes"] = nodes

	this.Show()
}

func (this *UpgradeRemoteAction) RunPost(params struct {
	NodeId int64

	Must *actions.Must
}) {
	_, err := this.RPC().NodeRPC().UpgradeNode(this.AdminContext(), &pb.UpgradeNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLogInfo(codes.Node_LogUpgradeNodeRemotely, params.NodeId)

	this.Success()
}
