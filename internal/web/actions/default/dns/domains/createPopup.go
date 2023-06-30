package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"strings"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	ProviderId int64
}) {
	this.Data["providerId"] = params.ProviderId

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	ProviderId int64
	Name       string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// TODO 检查ProviderId

	params.Must.
		Field("name", params.Name).
		Require("请输入域名")

	// 校验域名
	domain := strings.ToLower(params.Name)
	domain = strings.Replace(domain, " ", "", -1)
	if !domainutils.ValidateDomainFormat(domain) {
		this.Fail("域名格式不正确，请修改后重新提交")
	}

	createResp, err := this.RPC().DNSDomainRPC().CreateDNSDomain(this.AdminContext(), &pb.CreateDNSDomainRequest{
		DnsProviderId: params.ProviderId,
		Name:          domain,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	defer this.CreateLogInfo(codes.DNS_LogCreateDomain, createResp.DnsDomainId)

	this.Success()
}
