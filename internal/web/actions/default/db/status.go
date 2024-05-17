// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package db

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type StatusAction struct {
	actionutils.ParentAction
}

func (this *StatusAction) RunPost(params struct {
	NodeId int64
}) {
	statusResp, err := this.RPC().DBNodeRPC().CheckDBNodeStatus(this.AdminContext(), &pb.CheckDBNodeStatusRequest{DbNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var status = statusResp.DbNodeStatus
	if status != nil {
		this.Data["status"] = maps.Map{
			"isOk":    status.IsOk,
			"error":   status.Error,
			"size":    numberutils.FormatBytes(status.Size),
			"version": status.Version,
		}
	} else {
		this.Data["status"] = nil
	}

	this.Success()
}
