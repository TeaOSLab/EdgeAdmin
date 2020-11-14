package issues

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"net"
)

type UpdateNodePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateNodePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateNodePopupAction) RunGet(params struct {
	NodeId int64
}) {
	this.Data["nodeId"] = params.NodeId

	dnsInfoResp, err := this.RPC().NodeRPC().FindEnabledNodeDNS(this.AdminContext(), &pb.FindEnabledNodeDNSRequest{NodeId: params.NodeId})
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
	this.Data["route"] = dnsInfo.Route
	this.Data["domainId"] = dnsInfo.DnsDomainId

	// 读取所有线路
	routes := []string{}
	if dnsInfo.DnsDomainId > 0 {
		routesResp, err := this.RPC().DNSDomainRPC().FindAllDNSDomainRoutes(this.AdminContext(), &pb.FindAllDNSDomainRoutesRequest{DnsDomainId: dnsInfo.DnsDomainId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(routesResp.Routes) > 0 {
			routes = routesResp.Routes
		}
	}
	this.Data["routes"] = routes

	if len(routes) > 0 && !lists.ContainsString(routes, dnsInfo.Route) {
		this.Data["route"] = routes[0]
	}

	this.Show()
}

func (this *UpdateNodePopupAction) RunPost(params struct {
	NodeId   int64
	IpAddr   string
	DomainId int64
	Route    string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 操作日志
	this.CreateLog(oplogs.LevelInfo, "修改节点 %d 的DNS设置", params.NodeId)

	params.Must.
		Field("ipAddr", params.IpAddr).
		Require("请输入IP地址")

	if net.ParseIP(params.IpAddr) == nil {
		this.FailField("ipAddr", "请输入正确的IP地址")
	}

	// 执行修改
	_, err := this.RPC().NodeRPC().UpdateNodeDNS(this.AdminContext(), &pb.UpdateNodeDNSRequest{
		NodeId:      params.NodeId,
		IpAddr:      params.IpAddr,
		DnsDomainId: params.DomainId,
		Route:       params.Route,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
