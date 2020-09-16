package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
)

// 删除Header
type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	HeaderPolicyId int64
	Type           string
	HeaderId       int64
}) {
	policyConfigResp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{
		HeaderPolicyId: params.HeaderPolicyId,
	})
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

	switch params.Type {
	case "addHeader":
		result := []*shared.HTTPHeaderConfig{}
		for _, h := range policyConfig.AddHeaders {
			if h.Id != params.HeaderId {
				result = append(result, h)
			}
		}
		resultJSON, err := json.Marshal(result)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicyAddingHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicyAddingHeadersRequest{
			HeaderPolicyId: params.HeaderPolicyId,
			HeadersJSON:    resultJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	case "setHeader":
		result := []*shared.HTTPHeaderConfig{}
		for _, h := range policyConfig.SetHeaders {
			if h.Id != params.HeaderId {
				result = append(result, h)
			}
		}
		resultJSON, err := json.Marshal(result)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicySettingHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicySettingHeadersRequest{
			HeaderPolicyId: params.HeaderPolicyId,
			HeadersJSON:    resultJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	case "replace":
		result := []*shared.HTTPHeaderConfig{}
		for _, h := range policyConfig.ReplaceHeaders {
			if h.Id != params.HeaderId {
				result = append(result, h)
			}
		}
		resultJSON, err := json.Marshal(result)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicyReplacingHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicyReplacingHeadersRequest{
			HeaderPolicyId: params.HeaderPolicyId,
			HeadersJSON:    resultJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	case "addTrailer":
		result := []*shared.HTTPHeaderConfig{}
		for _, h := range policyConfig.AddTrailers {
			if h.Id != params.HeaderId {
				result = append(result, h)
			}
		}
		resultJSON, err := json.Marshal(result)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicyAddingTrailers(this.AdminContext(), &pb.UpdateHTTPHeaderPolicyAddingTrailersRequest{
			HeaderPolicyId: params.HeaderPolicyId,
			HeadersJSON:    resultJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Success()
}
