// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package dns

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "update")
	this.SecondMenu("dns")
}

func (this *IndexAction) RunGet(params struct {
	NodeId int64
}) {
	node, err := nodeutils.InitNodeInfo(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// DNS相关
	var clusters = []*pb.NodeCluster{node.NodeCluster}
	clusters = append(clusters, node.SecondaryNodeClusters...)
	var allDNSRouteMaps = map[int64][]maps.Map{} // domain id => routes
	var routeMaps = map[int64][]maps.Map{}       // domain id => routes
	var domainIds = []int64{}
	for _, cluster := range clusters {
		dnsInfoResp, err := this.RPC().NodeRPC().FindEnabledNodeDNS(this.AdminContext(), &pb.FindEnabledNodeDNSRequest{
			NodeId:        params.NodeId,
			NodeClusterId: cluster.Id,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var dnsInfo = dnsInfoResp.Node
		if dnsInfo.DnsDomainId <= 0 || len(dnsInfo.DnsDomainName) == 0 {
			continue
		}
		var domainId = dnsInfo.DnsDomainId

		// remove same domain
		if lists.ContainsInt64(domainIds, domainId) {
			continue
		}
		domainIds = append(domainIds, domainId)

		var domainName = dnsInfo.DnsDomainName
		if len(dnsInfo.Routes) > 0 {
			for _, route := range dnsInfo.Routes {
				routeMaps[domainId] = append(routeMaps[domainId], maps.Map{
					"domainId":   domainId,
					"domainName": domainName,
					"code":       route.Code,
					"name":       route.Name,
				})
			}
		}

		// 所有线路选项
		routesResp, err := this.RPC().DNSDomainRPC().FindAllDNSDomainRoutes(this.AdminContext(), &pb.FindAllDNSDomainRoutesRequest{DnsDomainId: dnsInfoResp.Node.DnsDomainId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, route := range routesResp.Routes {
			allDNSRouteMaps[domainId] = append(allDNSRouteMaps[domainId], maps.Map{
				"domainId":   domainId,
				"domainName": domainName,
				"name":       route.Name,
				"code":       route.Code,
			})
		}
	}

	var domainRoutes = []maps.Map{}
	for _, m := range routeMaps {
		domainRoutes = append(domainRoutes, m...)
	}
	this.Data["dnsRoutes"] = domainRoutes

	var allDomainRoutes = []maps.Map{}
	for _, m := range allDNSRouteMaps {
		allDomainRoutes = append(allDomainRoutes, m...)
	}
	this.Data["allDNSRoutes"] = allDomainRoutes

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	NodeId        int64
	DnsDomainId   int64
	DnsRoutesJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.NodeDNS_LogUpdateNodeDNS, params.NodeId)

	var rawRouteCodes = []string{}
	if len(params.DnsRoutesJSON) > 0 {
		err := json.Unmarshal(params.DnsRoutesJSON, &rawRouteCodes)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// remove duplications
	var dnsRouteCodes = []string{}
	for _, routeCode := range rawRouteCodes {
		if !lists.ContainsString(dnsRouteCodes, routeCode) {
			dnsRouteCodes = append(dnsRouteCodes, routeCode)
		}
	}

	_, err := this.RPC().NodeRPC().UpdateNodeDNS(this.AdminContext(), &pb.UpdateNodeDNSRequest{
		NodeId:      params.NodeId,
		IpAddr:      "",
		DnsDomainId: 0,
		Routes:      dnsRouteCodes,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
