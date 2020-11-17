package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("dns")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	dnsInfoResp, err := this.RPC().ServerRPC().FindEnabledServerDNS(this.AdminContext(), &pb.FindEnabledServerDNSRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["dnsName"] = dnsInfoResp.DnsName
	if dnsInfoResp.Domain != nil {
		this.Data["dnsDomain"] = dnsInfoResp.Domain.Name
	} else {
		this.Data["dnsDomain"] = ""
	}

	this.Show()
}
