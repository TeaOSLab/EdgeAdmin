package headers

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("header")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	webConfigResp, err := this.RPC().ServerRPC().FindAndInitServerWebConfig(this.AdminContext(), &pb.FindAndInitServerWebRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	webConfig := &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(webConfigResp.Config, webConfig)
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
		webConfigResp, err := this.RPC().ServerRPC().FindAndInitServerWebConfig(this.AdminContext(), &pb.FindAndInitServerWebRequest{ServerId: params.ServerId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		err = json.Unmarshal(webConfigResp.Config, webConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["requestHeaderPolicy"] = webConfig.RequestHeaders
	this.Data["responseHeaderPolicy"] = webConfig.ResponseHeaders

	this.Show()
}
