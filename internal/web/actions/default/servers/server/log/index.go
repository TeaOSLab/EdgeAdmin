package log

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "log", "")
	this.SecondMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ServerId  int64
	RequestId string
	Ip        string
	Domain    string
	Keyword   string
}) {
	this.Data["serverId"] = params.ServerId
	this.Data["requestId"] = params.RequestId
	this.Data["ip"] = params.Ip
	this.Data["domain"] = params.Domain
	this.Data["keyword"] = params.Keyword
	this.Data["path"] = this.Request.URL.Path

	// 记录最近使用
	_, err := this.RPC().LatestItemRPC().IncreaseLatestItem(this.AdminContext(), &pb.IncreaseLatestItemRequest{
		ItemType: "server",
		ItemId:   params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId  int64
	RequestId string
	Keyword   string
	Ip        string
	Domain    string

	Must *actions.Must
}) {
	isReverse := len(params.RequestId) > 0
	accessLogsResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
		ServerId:  params.ServerId,
		RequestId: params.RequestId,
		Size:      20,
		Day:       timeutil.Format("Ymd"),
		Keyword:   params.Keyword,
		Ip:        params.Ip,
		Domain:    params.Domain,
		Reverse:   isReverse,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	ipList := []string{}
	accessLogs := accessLogsResp.HttpAccessLogs
	if len(accessLogs) == 0 {
		accessLogs = []*pb.HTTPAccessLog{}
	} else {
		for _, accessLog := range accessLogs {
			if len(accessLog.RemoteAddr) > 0 {
				if !lists.ContainsString(ipList, accessLog.RemoteAddr) {
					ipList = append(ipList, accessLog.RemoteAddr)
				}
			}
		}
	}
	this.Data["accessLogs"] = accessLogs
	if len(accessLogs) > 0 {
		this.Data["requestId"] = accessLogs[0].RequestId
	} else {
		this.Data["requestId"] = params.RequestId
	}
	this.Data["hasMore"] = accessLogsResp.HasMore

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

	this.Success()
}
