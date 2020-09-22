package headers

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/webutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
}

func (this *IndexAction) RunGet(params struct {
	LocationId int64
}) {
	webConfig, err := webutils.FindWebConfigWithLocationId(this.Parent(), params.LocationId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 初始化Header
	webResp, err := this.RPC().HTTPWebRPC().FindEnabledHTTPWeb(this.AdminContext(), &pb.FindEnabledHTTPWebRequest{WebId: webConfig.Id})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	web := webResp.Web
	if web == nil {
		this.ErrorPage(errors.New("web should not be nil"))
		return
	}
	isChanged := false
	if web.RequestHeaderPolicyId <= 0 {
		createHeaderPolicyResp, err := this.RPC().HTTPHeaderPolicyRPC().CreateHTTPHeaderPolicy(this.AdminContext(), &pb.CreateHTTPHeaderPolicyRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		headerPolicyId := createHeaderPolicyResp.HeaderPolicyId
		_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebRequestHeaderPolicy(this.AdminContext(), &pb.UpdateHTTPWebRequestHeaderPolicyRequest{
			WebId:          web.Id,
			HeaderPolicyId: headerPolicyId,
		})
		isChanged = true
	}
	if web.ResponseHeaderPolicyId <= 0 {
		createHeaderPolicyResp, err := this.RPC().HTTPHeaderPolicyRPC().CreateHTTPHeaderPolicy(this.AdminContext(), &pb.CreateHTTPHeaderPolicyRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		headerPolicyId := createHeaderPolicyResp.HeaderPolicyId
		_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebResponseHeaderPolicy(this.AdminContext(), &pb.UpdateHTTPWebResponseHeaderPolicyRequest{
			WebId:          web.Id,
			HeaderPolicyId: headerPolicyId,
		})
		isChanged = true
	}

	// 重新获取配置
	if isChanged {
		webConfig, err = webutils.FindWebConfigWithLocationId(this.Parent(), params.LocationId)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["requestHeaderPolicy"] = webConfig.RequestHeaders
	this.Data["responseHeaderPolicy"] = webConfig.ResponseHeaders

	this.Show()
}
