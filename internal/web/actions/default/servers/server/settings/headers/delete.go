package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
)

// DeleteAction 删除Header
type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	HeaderPolicyId int64
	Type           string
	HeaderId       int64
}) {
	defer this.CreateLogInfo(codes.ServerHTTPHeader_LogDeleteHeader, params.HeaderPolicyId, params.HeaderId)

	policyConfigResp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{
		HttpHeaderPolicyId: params.HeaderPolicyId,
	})
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

	switch params.Type {
	case "setHeader":
		result := []*shared.HTTPHeaderRef{}
		for _, h := range policyConfig.SetHeaderRefs {
			if h.HeaderId != params.HeaderId {
				result = append(result, h)
			}
		}
		resultJSON, err := json.Marshal(result)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicySettingHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicySettingHeadersRequest{
			HttpHeaderPolicyId: params.HeaderPolicyId,
			HeadersJSON:        resultJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Success()
}
