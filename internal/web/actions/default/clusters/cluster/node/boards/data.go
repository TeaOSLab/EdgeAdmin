// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package boards

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type DataAction struct {
	actionutils.ParentAction
}

func (this *DataAction) RunPost(params struct {
	ClusterId int64
	NodeId    int64
}) {
	resp, err := this.RPC().ServerStatBoardRPC().ComposeServerStatNodeBoard(this.AdminContext(), &pb.ComposeServerStatNodeBoardRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["board"] = maps.Map{
		"isActive":            resp.IsActive,
		"trafficInBytes":      resp.TrafficInBytes,
		"trafficOutBytes":     resp.TrafficOutBytes,
		"countConnections":    resp.CountConnections,
		"countRequests":       resp.CountRequests,
		"countAttackRequests": resp.CountAttackRequests,
		"cpuUsage":            resp.CpuUsage,
		"memoryUsage":         resp.MemoryUsage,
		"memoryTotalSize":     resp.MemoryTotalSize,
		"load":                resp.Load,
		"cacheDiskSize":       resp.CacheDiskSize,
		"cacheMemorySize":     resp.CacheMemorySize,
	}

	this.Success()
}
