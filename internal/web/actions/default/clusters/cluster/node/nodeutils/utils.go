// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package nodeutils

import (
	"errors"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	executils "github.com/TeaOSLab/EdgeAdmin/internal/utils/exec"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"gopkg.in/yaml.v3"
	"os"
	"runtime"
	"strconv"
	"time"
)

// InitNodeInfo 初始化节点信息
func InitNodeInfo(parentAction *actionutils.ParentAction, nodeId int64) (*pb.Node, error) {
	// 节点信息（用于菜单）
	nodeResp, err := parentAction.RPC().NodeRPC().FindEnabledNode(parentAction.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: nodeId})
	if err != nil {
		return nil, err
	}
	if nodeResp.Node == nil {
		return nil, errors.New("node '" + strconv.FormatInt(nodeId, 10) + "' not found")
	}
	var node = nodeResp.Node

	info, err := parentAction.RPC().NodeRPC().FindEnabledNodeConfigInfo(parentAction.AdminContext(), &pb.FindEnabledNodeConfigInfoRequest{NodeId: nodeId})
	if err != nil {
		return nil, err
	}

	var groupMap maps.Map
	if node.NodeGroup != nil {
		groupMap = maps.Map{
			"id":   node.NodeGroup.Id,
			"name": node.NodeGroup.Name,
		}
	}

	parentAction.Data["node"] = maps.Map{
		"id":    node.Id,
		"name":  node.Name,
		"isOn":  node.IsOn,
		"isUp":  node.IsUp,
		"group": groupMap,
		"level": node.Level,
	}
	var clusterId int64 = 0
	if node.NodeCluster != nil {
		parentAction.Data["clusterId"] = node.NodeCluster.Id
		clusterId = node.NodeCluster.Id
	}

	// 左侧菜单
	var prefix = "/clusters/cluster/node"
	var query = "clusterId=" + types.String(clusterId) + "&nodeId=" + types.String(nodeId)
	var menuItem = parentAction.Data.GetString("secondMenuItem")

	var menuItems = []maps.Map{
		{
			"name":     parentAction.Lang(codes.NodeMenu_SettingBasic),
			"url":      prefix + "/update?" + query,
			"isActive": menuItem == "basic",
		},
		{
			"name":     parentAction.Lang(codes.NodeMenu_SettingDNS),
			"url":      prefix + "/settings/dns?" + query,
			"isActive": menuItem == "dns",
			"isOn":     info.HasDNSInfo,
		},
		{
			"name":     parentAction.Lang(codes.NodeMenu_SettingCache),
			"url":      prefix + "/settings/cache?" + query,
			"isActive": menuItem == "cache",
			"isOn":     info.HasCacheInfo,
		},
		{
			"name":     parentAction.Lang(codes.NodeMenu_SettingDDoSProtection),
			"url":      prefix + "/settings/ddos-protection?" + query,
			"isActive": menuItem == "ddosProtection",
			"isOn":     info.HasDDoSProtection,
		},
		{
			"name": "-",
			"url":  "",
		},
	}
	menuItems = filterMenuItems(menuItems, menuItem, prefix, query, info, parentAction.LangCode())
	menuItems = append(menuItems, []maps.Map{
		{
			"name":     parentAction.Lang(codes.NodeMenu_SettingSSH),
			"url":      prefix + "/settings/ssh?" + query,
			"isActive": menuItem == "ssh",
			"isOn":     info.HasSSH,
		},
		{
			"name":     parentAction.Lang(codes.NodeMenu_SettingSystem),
			"url":      prefix + "/settings/system?" + query,
			"isActive": menuItem == "system",
			"isOn":     info.HasSystemSettings,
		},
	}...)
	parentAction.Data["leftMenuItems"] = menuItems

	return nodeResp.Node, nil
}

// InstallLocalNode 安装本地节点
func InstallLocalNode() error {
	var targetDir = Tea.Root
	var nodeDir = targetDir + "/edge-node"
	var apiAddr = "http://127.0.0.1:8001" // 先固定

	// 查找节点安装文件
	var zipFile = Tea.Root + "/edge-api/deploy/edge-node-linux-" + runtime.GOARCH + "-v" + teaconst.Version /** 默认和管理系统一致 **/ + ".zip"
	{
		stat, err := os.Stat(zipFile)
		if err != nil {
			if os.IsNotExist(err) {
				return errors.New("installer file not found in '" + zipFile + "'")
			}
			return fmt.Errorf("open installer file failed: %w", err)
		}
		if stat.IsDir() {
			return errors.New("invalid installer file '" + zipFile + "'")
		}
	}

	// 解压节点
	var unzip = utils.NewUnzip(zipFile, targetDir)
	err := unzip.Run()
	if err != nil {
		return fmt.Errorf("unzip installer file failed: %w", err)
	}

	// 创建节点
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return fmt.Errorf("create rpc client failed: %w", err)
	}
	var ctx = rpcClient.Context(0)
	nodeClustersResp, err := rpcClient.NodeClusterRPC().ListEnabledNodeClusters(ctx, &pb.ListEnabledNodeClustersRequest{
		IdDesc: true,
		Offset: 0,
		Size:   1,
	})
	if err != nil {
		return err
	}
	if len(nodeClustersResp.NodeClusters) == 0 {
		return errors.New("no clusters yet, please create a cluster at least")
	}
	var clusterId = nodeClustersResp.NodeClusters[0].Id

	// 检查节点是否已生成
	countNodesResp, err := rpcClient.NodeRPC().CountAllEnabledNodesMatch(ctx, &pb.CountAllEnabledNodesMatchRequest{
		NodeClusterId: clusterId,
	})
	if err != nil {
		return err
	}
	if countNodesResp.Count > 0 {
		// 这里先不判断是否有本地节点，只要有节点，就不允许再次执行
		return errors.New("there are already nodes in the cluster")
	}

	createNodeResp, err := rpcClient.NodeRPC().CreateNode(ctx, &pb.CreateNodeRequest{
		Name:          "本地节点",
		NodeClusterId: clusterId,
		NodeLogin:     nil,
		NodeGroupId:   0,
		DnsRoutes:     nil,
		NodeRegionId:  0,
	})
	if err != nil {
		return err
	}
	var nodeId = createNodeResp.NodeId
	nodeResp, err := rpcClient.NodeRPC().FindEnabledNode(ctx, &pb.FindEnabledNodeRequest{NodeId: nodeId})
	if err != nil {
		return err
	}
	if nodeResp.Node == nil {
		return errors.New("could not find local node with created id '" + types.String(nodeId) + "'")
	}
	var node = nodeResp.Node

	// 生成节点配置
	var apiConfig = &configs.APIConfig{
		RPCEndpoints:     []string{apiAddr},
		RPCDisableUpdate: true,
		NodeId:           node.UniqueId,
		Secret:           node.Secret,
	}
	apiConfigYAML, err := yaml.Marshal(apiConfig)
	if err != nil {
		return fmt.Errorf("encode config failed: %w", err)
	}
	err = os.WriteFile(nodeDir+"/configs/api_node.yaml", apiConfigYAML, 0666)
	if err != nil {
		return fmt.Errorf("write config file failed: %w", err)
	}

	// 测试节点
	{
		var cmd = executils.NewTimeoutCmd(5*time.Second, nodeDir+"/bin/edge-node", "test")
		cmd.WithStdout()
		cmd.WithStderr()
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("node test failed: %w", err)
		}
	}

	// 启动节点
	{
		var cmd = executils.NewTimeoutCmd(5*time.Second, nodeDir+"/bin/edge-node", "start")
		cmd.WithStdout()
		cmd.WithStderr()
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("node start failed: %w", err)
		}
	}

	return nil
}
