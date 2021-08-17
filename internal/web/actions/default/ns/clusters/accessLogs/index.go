// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package logs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"net"
	"regexp"
	"strings"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct {
	RequestId string
	Keyword   string
	Day       string
}) {
	day := strings.ReplaceAll(params.Day, "-", "")
	if !regexp.MustCompile(`^\d{8}$`).MatchString(day) {
		day = timeutil.Format("Ymd")
	}

	this.Data["keyword"] = params.Keyword
	this.Data["day"] = day[:4] + "-" + day[4:6] + "-" + day[6:]
	this.Data["path"] = this.Request.URL.Path

	var size = int64(10)

	resp, err := this.RPC().NSAccessLogRPC().ListNSAccessLogs(this.AdminContext(), &pb.ListNSAccessLogsRequest{
		RequestId:  params.RequestId,
		NsNodeId:   0,
		NsDomainId: 0,
		NsRecordId: 0,
		Size:       size,
		Day:        day,
		Keyword:    params.Keyword,
		Reverse:    false,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	ipList := []string{}
	nodeIds := []int64{}
	domainIds := []int64{}
	if len(resp.NsAccessLogs) == 0 {
		this.Data["accessLogs"] = []interface{}{}
	} else {
		this.Data["accessLogs"] = resp.NsAccessLogs
		for _, accessLog := range resp.NsAccessLogs {
			// IP
			if len(accessLog.RemoteAddr) > 0 {
				// 去掉端口
				ip, _, err := net.SplitHostPort(accessLog.RemoteAddr)
				if err == nil {
					accessLog.RemoteAddr = ip
					if !lists.ContainsString(ipList, ip) {
						ipList = append(ipList, ip)
					}
				}
			}

			// 节点
			if !lists.ContainsInt64(nodeIds, accessLog.NsNodeId) {
				nodeIds = append(nodeIds, accessLog.NsNodeId)
			}

			// 域名
			if !lists.ContainsInt64(domainIds, accessLog.NsDomainId) {
				domainIds = append(domainIds, accessLog.NsDomainId)
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
		prevResp, err := this.RPC().NSAccessLogRPC().ListNSAccessLogs(this.AdminContext(), &pb.ListNSAccessLogsRequest{
			RequestId:  params.RequestId,
			NsNodeId:   0,
			NsDomainId: 0,
			NsRecordId: 0,
			Day:        day,
			Keyword:    params.Keyword,
			Size:       size,
			Reverse:    true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if int64(len(prevResp.NsAccessLogs)) == size {
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
				if len(region.Isp) > 0 {
					region.Summary += " | " + region.Isp
				}
				regionMap[ip] = region.Summary
			}
		}
	}
	this.Data["regions"] = regionMap

	// 节点信息
	nodeMap := map[int64]interface{}{} // node id => { ... }
	for _, nodeId := range nodeIds {
		nodeResp, err := this.RPC().NSNodeRPC().FindEnabledNSNode(this.AdminContext(), &pb.FindEnabledNSNodeRequest{NsNodeId: nodeId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		node := nodeResp.NsNode
		if node != nil {
			nodeMap[node.Id] = maps.Map{
				"id":   node.Id,
				"name": node.Name,
				"cluster": maps.Map{
					"id":   node.NsCluster.Id,
					"name": node.NsCluster.Name,
				},
			}
		}
	}
	this.Data["nodes"] = nodeMap

	// 域名信息
	domainMap := map[int64]interface{}{} // domain id => { ... }
	for _, domainId := range domainIds {
		domainResp, err := this.RPC().NSDomainRPC().FindEnabledNSDomain(this.AdminContext(), &pb.FindEnabledNSDomainRequest{NsDomainId: domainId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		domain := domainResp.NsDomain
		if domain != nil {
			domainMap[domain.Id] = maps.Map{
				"id":   domain.Id,
				"name": domain.Name,
			}
		}
	}
	this.Data["domains"] = domainMap

	this.Show()
}
