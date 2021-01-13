package log

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type TodayAction struct {
	actionutils.ParentAction
}

func (this *TodayAction) Init() {
	this.Nav("", "log", "")
	this.SecondMenu("today")
}

func (this *TodayAction) RunGet(params struct {
	RequestId string
	ServerId  int64
	HasError  int
}) {
	size := int64(10)

	this.Data["path"] = this.Request.URL.Path
	this.Data["hasError"] = params.HasError

	resp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
		RequestId: params.RequestId,
		ServerId:  params.ServerId,
		HasError:  params.HasError > 0,
		Day:       timeutil.Format("Ymd"),
		Size:      size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	ipList := []string{}
	if len(resp.AccessLogs) == 0 {
		this.Data["accessLogs"] = []interface{}{}
	} else {
		this.Data["accessLogs"] = resp.AccessLogs
		for _, accessLog := range resp.AccessLogs {
			if len(accessLog.RemoteAddr) > 0 {
				if !lists.ContainsString(ipList, accessLog.RemoteAddr) {
					ipList = append(ipList, accessLog.RemoteAddr)
				}
			}
		}
	}
	this.Data["hasMore"] = resp.HasMore
	this.Data["nextRequestId"] = resp.RequestId

	// 上一个requestId
	this.Data["hasPrev"] = false
	this.Data["lastRequestId"] = ""
	if len(params.RequestId) > 0 {
		this.Data["hasPrev"] = true
		prevResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
			RequestId: params.RequestId,
			ServerId:  params.ServerId,
			HasError:  params.HasError > 0,
			Day:       timeutil.Format("Ymd"),
			Size:      size,
			Reverse:   true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if int64(len(prevResp.AccessLogs)) == size {
			this.Data["lastRequestId"] = prevResp.RequestId
		}
	}

	// 根据IP查询区域
	regionMap := map[string]string{} // ip => region
	if len(ipList) > 0 {
		resp, err := this.RPC().IPLibraryRPC().LookupIPRegions(this.AdminContext(), &pb.LookupIPRegionsRequest{IpList: ipList})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if resp.IpRegionMap != nil {
			for ip, region := range resp.IpRegionMap {
				regionMap[ip] = region.Summary
			}
		}
	}
	this.Data["regions"] = regionMap

	this.Show()
}
