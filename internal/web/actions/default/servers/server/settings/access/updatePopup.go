// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package access

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	PolicyId int64
}) {
	this.Data["authTypes"] = serverconfigs.FindAllHTTPAuthTypes()

	policyResp, err := this.RPC().HTTPAuthPolicyRPC().FindEnabledHTTPAuthPolicy(this.AdminContext(), &pb.FindEnabledHTTPAuthPolicyRequest{HttpAuthPolicyId: params.PolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policy := policyResp.HttpAuthPolicy
	if policy == nil {
		this.NotFound("httpAuthPolicy", params.PolicyId)
		return
	}

	var authParams = map[string]interface{}{}
	if len(policy.ParamsJSON) > 0 {
		err = json.Unmarshal(policy.ParamsJSON, &authParams)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["policy"] = maps.Map{
		"id":     policy.Id,
		"isOn":   policy.IsOn,
		"name":   policy.Name,
		"type":   policy.Type,
		"params": authParams,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	PolicyId int64

	Name string
	IsOn bool

	// BasicAuth
	HttpAuthBasicAuthUsersJSON []byte
	BasicAuthRealm             string
	BasicAuthCharset           string

	// SubRequest
	SubRequestURL           string
	SubRequestMethod        string
	SubRequestFollowRequest bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改HTTP认证 %d", params.PolicyId)

	policyResp, err := this.RPC().HTTPAuthPolicyRPC().FindEnabledHTTPAuthPolicy(this.AdminContext(), &pb.FindEnabledHTTPAuthPolicyRequest{HttpAuthPolicyId: params.PolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policy := policyResp.HttpAuthPolicy
	if policy == nil {
		this.NotFound("httpAuthPolicy", params.PolicyId)
		return
	}
	policyType := policy.Type

	params.Must.
		Field("name", params.Name).
		Require("请输入名称")

	var ref = &serverconfigs.HTTPAuthPolicyRef{IsOn: true}
	var paramsJSON []byte

	switch policyType {
	case serverconfigs.HTTPAuthTypeBasicAuth:
		users := []*serverconfigs.HTTPAuthBasicMethodUser{}
		err := json.Unmarshal(params.HttpAuthBasicAuthUsersJSON, &users)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(users) == 0 {
			this.Fail("请添加至少一个用户")
		}
		method := &serverconfigs.HTTPAuthBasicMethod{
			Users:   users,
			Realm:   params.BasicAuthRealm,
			Charset: params.BasicAuthCharset,
		}
		methodJSON, err := json.Marshal(method)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		paramsJSON = methodJSON
	case serverconfigs.HTTPAuthTypeSubRequest:
		params.Must.Field("subRequestURL", params.SubRequestURL).
			Require("请输入子请求URL")
		if params.SubRequestFollowRequest {
			params.SubRequestMethod = ""
		}
		method := &serverconfigs.HTTPAuthSubRequestMethod{
			URL:    params.SubRequestURL,
			Method: params.SubRequestMethod,
		}
		methodJSON, err := json.Marshal(method)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		paramsJSON = methodJSON
	default:
		this.Fail("不支持的认证类型'" + policyType + "'")
	}

	var paramsMap map[string]interface{}
	err = json.Unmarshal(paramsJSON, &paramsMap)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPAuthPolicyRPC().UpdateHTTPAuthPolicy(this.AdminContext(), &pb.UpdateHTTPAuthPolicyRequest{
		HttpAuthPolicyId: params.PolicyId,
		Name:             params.Name,
		ParamsJSON:       paramsJSON,
		IsOn:             params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	ref.AuthPolicy = &serverconfigs.HTTPAuthPolicy{
		Id:     params.PolicyId,
		Name:   params.Name,
		IsOn:   params.IsOn,
		Type:   policyType,
		Params: paramsMap,
	}

	this.Data["policyRef"] = ref
	this.Success()
}
