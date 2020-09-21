package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
)

type DeleteDeletingHeaderAction struct {
	actionutils.ParentAction
}

func (this *DeleteDeletingHeaderAction) RunPost(params struct {
	HeaderPolicyId int64
	HeaderName     string
}) {
	policyConfigResp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{HeaderPolicyId: params.HeaderPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyConfigJSON := policyConfigResp.HeaderPolicyJSON
	policyConfig := &shared.HTTPHeaderPolicy{}
	err = json.Unmarshal(policyConfigJSON, policyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	headerNames := []string{}
	for _, h := range policyConfig.DeletedHeaders {
		if h == params.HeaderName {
			continue
		}
		headerNames = append(headerNames, h)
	}
	_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicyDeletingHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicyDeletingHeadersRequest{
		HeaderPolicyId: params.HeaderPolicyId,
		HeaderNames:    headerNames,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
