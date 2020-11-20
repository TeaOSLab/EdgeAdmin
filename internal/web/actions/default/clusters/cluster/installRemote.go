package cluster

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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
	this.Data["leftMenuItems"] = LeftMenuItemsForInstall(params.ClusterId, "install")

	nodesResp, err := this.RPC().NodeRPC().FindAllNotInstalledNodesWithClusterId(this.AdminContext(), &pb.FindAllNotInstalledNodesWithClusterIdRequest{ClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	nodeMaps := []maps.Map{}
	for _, node := range nodesResp.Nodes {
		loginParams := maps.Map{}
		if node.Login != nil && len(node.Login.Params) > 0 {
			err := json.Unmarshal(node.Login.Params, &loginParams)
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
			"login":         node.Login,
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
	defer this.CreateLog(oplogs.LevelInfo, "远程安装节点 %d", params.NodeId)

	this.Success()
}
