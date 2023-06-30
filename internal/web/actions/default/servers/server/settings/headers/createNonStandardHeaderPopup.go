package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type CreateNonStandardPopupAction struct {
	actionutils.ParentAction
}

func (this *CreateNonStandardPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateNonStandardPopupAction) RunGet(params struct {
	HeaderPolicyId int64
	Type           string
}) {
	this.Data["headerPolicyId"] = params.HeaderPolicyId
	this.Data["type"] = params.Type

	this.Show()
}

func (this *CreateNonStandardPopupAction) RunPost(params struct {
	HeaderPolicyId int64
	Name           string

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.ServerHTTPHeader_LogCreateNonStandardHeader, params.HeaderPolicyId, params.Name)

	params.Must.
		Field("name", params.Name).
		Require("名称不能为空")

	policyConfigResp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{HttpHeaderPolicyId: params.HeaderPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var policyConfig = &shared.HTTPHeaderPolicy{}
	err = json.Unmarshal(policyConfigResp.HttpHeaderPolicyJSON, policyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var nonStandardHeaders = policyConfig.NonStandardHeaders
	nonStandardHeaders = append(nonStandardHeaders, params.Name)
	_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicyNonStandardHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicyNonStandardHeadersRequest{
		HttpHeaderPolicyId: params.HeaderPolicyId,
		HeaderNames:        nonStandardHeaders,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
