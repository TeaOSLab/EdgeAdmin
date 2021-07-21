// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipbox

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct {
	Ip string
}) {
	this.Data["ip"] = params.Ip

	accessLogsResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
		Day:     timeutil.Format("Ymd"),
		Keyword: "ip:" + params.Ip,
		Size:    10,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var accessLogs = accessLogsResp.HttpAccessLogs
	if len(accessLogs) == 0 {
		// 查询昨天
		accessLogsResp, err = this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
			Day:     timeutil.Format("Ymd", time.Now().AddDate(0, 0, -1)),
			Keyword: "ip:" + params.Ip,
			Size:    10,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		accessLogs = accessLogsResp.HttpAccessLogs

		if len(accessLogs) == 0 {
			accessLogs = []*pb.HTTPAccessLog{}
		}
	}
	this.Data["accessLogs"] = accessLogs

	this.Show()
}
