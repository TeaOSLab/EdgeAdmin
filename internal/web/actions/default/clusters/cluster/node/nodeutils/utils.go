// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package nodeutils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"strconv"
)

// InitNodeInfo 初始化节点信息
func InitNodeInfo(action actionutils.ActionInterface, nodeId int64) error {
	// 节点信息（用于菜单）
	nodeResp, err := action.RPC().NodeRPC().FindEnabledNode(action.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: nodeId})
	if err != nil {
		return err
	}
	if nodeResp.Node == nil {
		return errors.New("node '" + strconv.FormatInt(nodeId, 10) + "' not found")
	}
	var node = nodeResp.Node
	action.ViewData()["node"] = maps.Map{
		"id":   node.Id,
		"name": node.Name,
	}
	if node.NodeCluster != nil {
		action.ViewData()["clusterId"] = node.NodeCluster.Id
	}
	return nil
}
