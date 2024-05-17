// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package dbnodeutils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

// InitNode 初始化指标信息
func InitNode(parent *actionutils.ParentAction, nodeId int64) (*pb.DBNode, error) {
	client, err := rpc.SharedRPC()
	if err != nil {
		return nil, err
	}

	resp, err := client.DBNodeRPC().FindEnabledDBNode(parent.AdminContext(), &pb.FindEnabledDBNodeRequest{DbNodeId: nodeId})
	if err != nil {
		return nil, err
	}
	var node = resp.DbNode
	if node == nil {
		return nil, errors.New("not found db node with id '" + types.String(nodeId) + "'")
	}
	parent.Data["node"] = maps.Map{
		"id":   node.Id,
		"name": node.Name,
	}
	return node, nil
}
