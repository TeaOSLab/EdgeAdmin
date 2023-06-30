package cluster

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type InstallRemoteAction struct {
	actionutils.ParentAction
}

func (this *InstallRemoteAction) Init() {
	this.Nav("", "node", "install")
	this.SecondMenu("nodes")
}

func (this *InstallRemoteAction) RunGet(params struct {
	ClusterId int64
}) {
	this.Data["leftMenuItems"] = LeftMenuItemsForInstall(this.AdminContext(), params.ClusterId, "install", this.LangCode())

	nodesResp, err := this.RPC().NodeRPC().FindAllNotInstalledNodesWithNodeClusterId(this.AdminContext(), &pb.FindAllNotInstalledNodesWithNodeClusterIdRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	nodeMaps := []maps.Map{}
	for _, node := range nodesResp.Nodes {
		loginParams := maps.Map{}
		if node.NodeLogin != nil && len(node.NodeLogin.Params) > 0 {
			err := json.Unmarshal(node.NodeLogin.Params, &loginParams)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		installStatus := maps.Map{
			"isRunning":  false,
			"isFinished": false,
		}
		if node.InstallStatus != nil {
			installStatus = maps.Map{
				"isRunning":  node.InstallStatus.IsRunning,
				"isFinished": node.InstallStatus.IsFinished,
				"isOk":       node.InstallStatus.IsOk,
				"error":      node.InstallStatus.Error,
			}
		}

		nodeMaps = append(nodeMaps, maps.Map{
			"id":            node.Id,
			"isOn":          node.IsOn,
			"name":          node.Name,
			"addresses":     node.IpAddresses,
			"login":         node.NodeLogin,
			"loginParams":   loginParams,
			"installStatus": installStatus,
		})
	}
	this.Data["nodes"] = nodeMaps

	this.Show()
}

func (this *InstallRemoteAction) RunPost(params struct {
	NodeId int64

	Must *actions.Must
}) {
	_, err := this.RPC().NodeRPC().InstallNode(this.AdminContext(), &pb.InstallNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLogInfo(codes.Node_LogInstallNodeRemotely, params.NodeId)

	this.Success()
}
