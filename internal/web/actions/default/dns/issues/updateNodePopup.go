package issues

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net"
)

type UpdateNodePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateNodePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateNodePopupAction) RunGet(params struct {
	ClusterId int64
	NodeId    int64
	IpAddrId  int64
}) {
	this.Data["nodeId"] = params.NodeId

	dnsInfoResp, err := this.RPC().NodeRPC().FindEnabledNodeDNS(this.AdminContext(), &pb.FindEnabledNodeDNSRequest{
		NodeId:        params.NodeId,
		NodeClusterId: params.ClusterId,
		NodeIPAddrId:  params.IpAddrId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	dnsInfo := dnsInfoResp.Node
	if dnsInfo == nil {
		this.NotFound("node", params.NodeId)
		return
	}
	this.Data["ipAddr"] = dnsInfo.IpAddr
	this.Data["ipAddrId"] = dnsInfo.NodeIPAddressId
	this.Data["routes"] = domainutils.ConvertRoutesToMaps(dnsInfo)
	this.Data["domainId"] = dnsInfo.DnsDomainId
	this.Data["domainName"] = dnsInfo.DnsDomainName

	// 读取所有线路
	var allRouteMaps = []maps.Map{}
	if dnsInfo.DnsDomainId > 0 {
		routesResp, err := this.RPC().DNSDomainRPC().FindAllDNSDomainRoutes(this.AdminContext(), &pb.FindAllDNSDomainRoutesRequest{DnsDomainId: dnsInfo.DnsDomainId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(routesResp.Routes) > 0 {
			for _, route := range routesResp.Routes {
				allRouteMaps = append(allRouteMaps, maps.Map{
					"name":       route.Name,
					"code":       route.Code,
					"domainName": dnsInfo.DnsDomainName,
					"domainId":   dnsInfo.DnsDomainId,
				})
			}

			// 筛选
			var routes = domainutils.FilterRoutes(dnsInfo.Routes, routesResp.Routes)
			dnsInfo.Routes = routes
			this.Data["routes"] = domainutils.ConvertRoutesToMaps(dnsInfo)
		}
	}
	this.Data["allRoutes"] = allRouteMaps

	this.Show()
}

func (this *UpdateNodePopupAction) RunPost(params struct {
	NodeId        int64
	IpAddr        string
	IpAddrId      int64
	DomainId      int64
	DnsRoutesJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 操作日志
	defer this.CreateLogInfo(codes.DNS_LogUpdateNodeDNS, params.NodeId)

	var routes = []string{}
	if len(params.DnsRoutesJSON) > 0 {
		err := json.Unmarshal(params.DnsRoutesJSON, &routes)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	params.Must.
		Field("ipAddr", params.IpAddr).
		Require("请输入IP地址")

	if net.ParseIP(params.IpAddr) == nil {
		this.FailField("ipAddr", "请输入正确的IP地址")
	}

	// 执行修改
	_, err := this.RPC().NodeRPC().UpdateNodeDNS(this.AdminContext(), &pb.UpdateNodeDNSRequest{
		NodeId:          params.NodeId,
		IpAddr:          params.IpAddr,
		NodeIPAddressId: params.IpAddrId,
		DnsDomainId:     params.DomainId,
		Routes:          routes,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
