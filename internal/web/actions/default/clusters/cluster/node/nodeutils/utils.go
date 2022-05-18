// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package nodeutils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strconv"
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
			"name":     "基础设置",
			"url":      prefix + "/update?" + query,
			"isActive": menuItem == "basic",
		},
		{
			"name":     "DNS设置",
			"url":      prefix + "/settings/dns?" + query,
			"isActive": menuItem == "dns",
			"isOn":     info.HasDNSInfo,
		},
		{
			"name":     "缓存设置",
			"url":      prefix + "/settings/cache?" + query,
			"isActive": menuItem == "cache",
			"isOn":     info.HasCacheInfo,
		},
		{
			"name":     "DDOS防护",
			"url":      prefix + "/settings/ddos-protection?" + query,
			"isActive": menuItem == "ddosProtection",
			"isOn":     info.HasDDoSProtection,
		},
		{
			"name": "-",
			"url":  "",
		},
	}
	menuItems = filterMenuItems(menuItems, menuItem, prefix, query, info)
	menuItems = append(menuItems, []maps.Map{
		{
			"name":     "SSH设置",
			"url":      prefix + "/settings/ssh?" + query,
			"isActive": menuItem == "ssh",
			"isOn":     info.HasSSH,
		},
		{
			"name":     "系统设置",
			"url":      prefix + "/settings/system?" + query,
			"isActive": menuItem == "system",
			"isOn":     info.HasSystemSettings,
		},
	}...)
	parentAction.Data["leftMenuItems"] = menuItems

	return nodeResp.Node, nil
}
