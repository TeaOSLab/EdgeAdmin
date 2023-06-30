package domains

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	DomainId int64
}) {
	domainResp, err := this.RPC().DNSDomainRPC().FindDNSDomain(this.AdminContext(), &pb.FindDNSDomainRequest{DnsDomainId: params.DomainId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domain := domainResp.DnsDomain
	if domain == nil {
		this.NotFound("dnsDomain", params.DomainId)
		return
	}

	this.Data["domain"] = maps.Map{
		"id":   domain.Id,
		"name": domain.Name,
		"isOn": domain.IsOn,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	DomainId int64
	Name     string
	IsOn     bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// TODO 检查DomainId

	// 记录日志
	defer this.CreateLogInfo(codes.DNS_LogUpdateDomain, params.DomainId)

	params.Must.
		Field("name", params.Name).
		Require("请输入域名")

	// 校验域名
	domain := strings.ToLower(params.Name)
	domain = strings.Replace(domain, " ", "", -1)
	if !domainutils.ValidateDomainFormat(domain) {
		this.Fail("域名格式不正确，请修改后重新提交")
	}

	_, err := this.RPC().DNSDomainRPC().UpdateDNSDomain(this.AdminContext(), &pb.UpdateDNSDomainRequest{
		DnsDomainId: params.DomainId,
		Name:        domain,
		IsOn:        params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
