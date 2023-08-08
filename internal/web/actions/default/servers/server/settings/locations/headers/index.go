package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
}

func (this *IndexAction) RunGet(params struct {
	LocationId int64
}) {
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithLocationId(this.AdminContext(), params.LocationId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	webId := webConfig.Id
	this.Data["webId"] = webId

	isChanged := false

	if webConfig.RequestHeaderPolicy == nil {
		createHeaderPolicyResp, err := this.RPC().HTTPHeaderPolicyRPC().CreateHTTPHeaderPolicy(this.AdminContext(), &pb.CreateHTTPHeaderPolicyRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		headerPolicyId := createHeaderPolicyResp.HttpHeaderPolicyId
		ref := &shared.HTTPHeaderPolicyRef{
			IsPrior:        false,
			IsOn:           true,
			HeaderPolicyId: headerPolicyId,
		}
		refJSON, err := json.Marshal(ref)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebRequestHeader(this.AdminContext(), &pb.UpdateHTTPWebRequestHeaderRequest{
			HttpWebId:  webId,
			HeaderJSON: refJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		isChanged = true
	}
	if webConfig.ResponseHeaderPolicy == nil {
		createHeaderPolicyResp, err := this.RPC().HTTPHeaderPolicyRPC().CreateHTTPHeaderPolicy(this.AdminContext(), &pb.CreateHTTPHeaderPolicyRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		headerPolicyId := createHeaderPolicyResp.HttpHeaderPolicyId
		ref := &shared.HTTPHeaderPolicyRef{
			IsPrior:        false,
			IsOn:           true,
			HeaderPolicyId: headerPolicyId,
		}
		refJSON, err := json.Marshal(ref)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebResponseHeader(this.AdminContext(), &pb.UpdateHTTPWebResponseHeaderRequest{
			HttpWebId:  webId,
			HeaderJSON: refJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		isChanged = true
	}

	// 重新获取配置
	if isChanged {
		webConfig, err = dao.SharedHTTPWebDAO.FindWebConfigWithLocationId(this.AdminContext(), params.LocationId)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["requestHeaderRef"] = webConfig.RequestHeaderPolicyRef
	this.Data["requestHeaderPolicy"] = webConfig.RequestHeaderPolicy
	this.Data["responseHeaderRef"] = webConfig.ResponseHeaderPolicyRef
	this.Data["responseHeaderPolicy"] = webConfig.ResponseHeaderPolicy

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId              int64
	Type               string
	RequestHeaderJSON  []byte
	ResponseHeaderJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.ServerHTTPHeader_LogUpdateHTTPHeaders, params.WebId)

	// TODO 检查配置

	switch params.Type {
	case "request":
		_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebRequestHeader(this.AdminContext(), &pb.UpdateHTTPWebRequestHeaderRequest{
			HttpWebId:  params.WebId,
			HeaderJSON: params.RequestHeaderJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	case "response":
		_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebResponseHeader(this.AdminContext(), &pb.UpdateHTTPWebResponseHeaderRequest{
			HttpWebId:  params.WebId,
			HeaderJSON: params.ResponseHeaderJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Success()
}
