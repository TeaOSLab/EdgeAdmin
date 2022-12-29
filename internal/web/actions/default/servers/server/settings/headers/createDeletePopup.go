package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type CreateDeletePopupAction struct {
	actionutils.ParentAction
}

func (this *CreateDeletePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateDeletePopupAction) RunGet(params struct {
	HeaderPolicyId int64
	Type           string
}) {
	this.Data["headerPolicyId"] = params.HeaderPolicyId
	this.Data["type"] = params.Type

	this.Show()
}

func (this *CreateDeletePopupAction) RunPost(params struct {
	HeaderPolicyId int64
	Name           string

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "添加删除的Header HeaderPolicyId: %d, Name: %s", params.HeaderPolicyId, params.Name)

	params.Must.
		Field("name", params.Name).
		Require("名称不能为空")

	policyConfigResp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{HttpHeaderPolicyId: params.HeaderPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyConfig := &shared.HTTPHeaderPolicy{}
	err = json.Unmarshal(policyConfigResp.HttpHeaderPolicyJSON, policyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	deleteHeaders := policyConfig.DeleteHeaders
	deleteHeaders = append(deleteHeaders, params.Name)
	_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicyDeletingHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicyDeletingHeadersRequest{
		HttpHeaderPolicyId: params.HeaderPolicyId,
		HeaderNames:        deleteHeaders,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
