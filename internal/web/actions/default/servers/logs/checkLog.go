// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package logs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

// CheckLogAction 检查是否有日志
type CheckLogAction struct {
	actionutils.ParentAction
}

func (this *CheckLogAction) RunPost(params struct {
	Day       string
	Partition int32

	ServerId          int64
	RequestId         string
	ClusterId         int64
	NodeId            int64
	HasError          bool
	HasFirewallPolicy bool
	Keyword           string
	Ip                string
	Domain            string

	Hour string
}) {
	resp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
		Partition:         params.Partition,
		RequestId:         params.RequestId,
		NodeClusterId:     params.ClusterId,
		NodeId:            params.NodeId,
		ServerId:          params.ServerId,
		HasError:          params.HasError,
		HasFirewallPolicy: params.HasFirewallPolicy,
		Day:               params.Day,
		HourFrom:          params.Hour,
		HourTo:            params.Hour,
		Keyword:           params.Keyword,
		Ip:                params.Ip,
		Domain:            params.Domain,
		Size:              1,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["hasLogs"] = len(resp.HttpAccessLogs) > 0

	this.Success()
}
