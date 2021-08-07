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
	Keyword   string
	Ip        string
	Domain    string
}) {
	size := int64(10)

	this.Data["path"] = this.Request.URL.Path
	this.Data["hasError"] = params.HasError
	this.Data["keyword"] = params.Keyword
	this.Data["ip"] = params.Ip
	this.Data["domain"] = params.Domain

	resp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
		RequestId: params.RequestId,
		ServerId:  params.ServerId,
		HasError:  params.HasError > 0,
		Day:       timeutil.Format("Ymd"),
		Keyword:   params.Keyword,
		Ip:        params.Ip,
		Domain:    params.Domain,
		Size:      size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	ipList := []string{}
	if len(resp.HttpAccessLogs) == 0 {
		this.Data["accessLogs"] = []interface{}{}
	} else {
		this.Data["accessLogs"] = resp.HttpAccessLogs
		for _, accessLog := range resp.HttpAccessLogs {
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
			Keyword:   params.Keyword,
			Ip:        params.Ip,
			Domain:    params.Domain,
			Size:      size,
			Reverse:   true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if int64(len(prevResp.HttpAccessLogs)) == size {
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
