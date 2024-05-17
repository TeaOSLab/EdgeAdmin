// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package logs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"regexp"
	"strings"
)

// HasLogsAction 检查某个分区是否有日志
type HasLogsAction struct {
	actionutils.ParentAction
}

func (this *HasLogsAction) RunPost(params struct {
	ClusterId int64
	NodeId    int64
	Day       string
	Hour      string
	Keyword   string
	Ip        string
	Domain    string
	HasError  int
	HasWAF    int
	Partition int32 `default:"-1"`

	RequestId string
	ServerId  int64
}) {
	if len(params.Day) == 0 {
		params.Day = timeutil.Format("Y-m-d")
	}

	var day = params.Day

	if len(day) > 0 && regexp.MustCompile(`\d{4}-\d{2}-\d{2}`).MatchString(day) {
		day = strings.ReplaceAll(day, "-", "")
	}

	resp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
		Partition:         params.Partition,
		RequestId:         params.RequestId,
		NodeClusterId:     params.ClusterId,
		NodeId:            params.NodeId,
		ServerId:          params.ServerId,
		HasError:          params.HasError > 0,
		HasFirewallPolicy: params.HasWAF > 0,
		Day:               day,
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
