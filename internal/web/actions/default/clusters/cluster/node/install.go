package node

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"path/filepath"
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

	// 最近运行目录
	var exeRoot = ""
	if len(node.StatusJSON) > 0 {
		var nodeStatus = &nodeconfigs.NodeStatus{}
		err = json.Unmarshal(node.StatusJSON, nodeStatus)
		if err == nil {
			var exePath = nodeStatus.ExePath
			if len(exePath) > 0 {
				exeRoot = filepath.Dir(filepath.Dir(exePath))
			}
		}
	}
	this.Data["exeRoot"] = exeRoot

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
	var apiNodes = apiNodesResp.ApiNodes
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

	// 安装文件
	var installerFiles = clusterutils.ListInstallerFiles()
	this.Data["installerFiles"] = installerFiles

	// SSH主机地址
	this.Data["sshAddr"] = ""
	if node.NodeLogin != nil && node.NodeLogin.Type == "ssh" && !utils.JSONIsNull(node.NodeLogin.Params) {
		var loginParams = maps.Map{}
		err = json.Unmarshal(node.NodeLogin.Params, &loginParams)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		var host = loginParams.GetString("host")
		if len(host) > 0 {
			var port = loginParams.GetString("port")
			if port == "0" {
				port = "22"
			}
			this.Data["sshAddr"] = configutils.QuoteIP(host) + ":" + port
		}
	}

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
	defer this.CreateLogInfo(codes.Node_LogInstallNode, params.NodeId)

	this.Success()
}
