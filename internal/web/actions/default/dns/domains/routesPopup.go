package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type RoutesPopupAction struct {
	actionutils.ParentAction
}

func (this *RoutesPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *RoutesPopupAction) RunGet(params struct {
	DomainId int64
}) {
	routesResp, err := this.RPC().DNSDomainRPC().FindAllDNSDomainRoutes(this.AdminContext(), &pb.FindAllDNSDomainRoutesRequest{DnsDomainId: params.DomainId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if len(routesResp.Routes) == 0 {
		routesResp.Routes = []string{}
	}
	this.Data["routes"] = routesResp.Routes

	this.Show()
}
