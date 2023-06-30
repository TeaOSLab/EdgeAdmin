package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/groups/group/servergrouputils"
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
	this.Nav("", "setting", "index")
	this.SecondMenu("header")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
}) {
	_, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "headers")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerGroupId(this.AdminContext(), params.GroupId)
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
		webConfig, err = dao.SharedHTTPWebDAO.FindWebConfigWithServerGroupId(this.AdminContext(), params.GroupId)
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

	switch params.Type {
	case "request":
		// 检查配置
		var requestHeaderRef = &shared.HTTPHeaderPolicyRef{}
		err := json.Unmarshal(params.RequestHeaderJSON, requestHeaderRef)
		if err != nil {
			this.Fail("请求Header配置校验失败：" + err.Error() + ", JSON: " + string(params.RequestHeaderJSON))
		}

		_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebRequestHeader(this.AdminContext(), &pb.UpdateHTTPWebRequestHeaderRequest{
			HttpWebId:  params.WebId,
			HeaderJSON: params.RequestHeaderJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	case "response":
		// 校验配置
		var responseHeaderRef = &shared.HTTPHeaderPolicyRef{}
		err := json.Unmarshal(params.ResponseHeaderJSON, responseHeaderRef)
		if err != nil {
			this.Fail("响应Header配置校验失败：" + err.Error() + ", JSON: " + string(params.ResponseHeaderJSON))
		}

		_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebResponseHeader(this.AdminContext(), &pb.UpdateHTTPWebResponseHeaderRequest{
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
