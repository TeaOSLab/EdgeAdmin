package headers

import (
	"encoding/json"
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
}) {
	this.Data["headerPolicyId"] = params.HeaderPolicyId

	this.Show()
}

func (this *CreateDeletePopupAction) RunPost(params struct {
	HeaderPolicyId int64
	Name           string

	Must *actions.Must
}) {
	policyConfigResp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{HeaderPolicyId: params.HeaderPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyConfig := &shared.HTTPHeaderPolicy{}
	err = json.Unmarshal(policyConfigResp.Config, policyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	deleteHeaders := policyConfig.DeletedHeaders
	deleteHeaders = append(deleteHeaders, params.Name)
	_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicyDeletingHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicyDeletingHeadersRequest{
		HeaderPolicyId: params.HeaderPolicyId,
		HeaderNames:    deleteHeaders,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
