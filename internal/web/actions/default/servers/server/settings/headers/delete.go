package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
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
	defer this.CreateLog(oplogs.LevelInfo, "删除请求Header，HeaderPolicyId:%d, HeaderId:%d", params.HeaderPolicyId, params.HeaderId)

	policyConfigResp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{
		HeaderPolicyId: params.HeaderPolicyId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyConfig := &shared.HTTPHeaderPolicy{}
	err = json.Unmarshal(policyConfigResp.HeaderPolicyJSON, policyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	switch params.Type {
	case "addHeader":
		result := []*shared.HTTPHeaderRef{}
		for _, h := range policyConfig.AddHeaderRefs {
			if h.HeaderId != params.HeaderId {
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
			HeaderPolicyId: params.HeaderPolicyId,
			HeadersJSON:    resultJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	case "replace":
		result := []*shared.HTTPHeaderRef{}
		for _, h := range policyConfig.ReplaceHeaderRefs {
			if h.HeaderId != params.HeaderId {
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
		result := []*shared.HTTPHeaderRef{}
		for _, h := range policyConfig.AddTrailerRefs {
			if h.HeaderId != params.HeaderId {
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
