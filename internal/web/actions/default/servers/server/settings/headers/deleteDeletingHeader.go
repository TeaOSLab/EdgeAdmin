package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
	// 日志
	defer this.CreateLogInfo(codes.ServerHTTPHeader_LogDeleteDeletingHeader, params.HeaderPolicyId, params.HeaderName)

	policyConfigResp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{HttpHeaderPolicyId: params.HeaderPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var policyConfigJSON = policyConfigResp.HttpHeaderPolicyJSON
	var policyConfig = &shared.HTTPHeaderPolicy{}
	err = json.Unmarshal(policyConfigJSON, policyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var headerNames = []string{}
	for _, h := range policyConfig.DeleteHeaders {
		if h == params.HeaderName {
			continue
		}
		headerNames = append(headerNames, h)
	}
	_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicyDeletingHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicyDeletingHeadersRequest{
		HttpHeaderPolicyId: params.HeaderPolicyId,
		HeaderNames:        headerNames,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
