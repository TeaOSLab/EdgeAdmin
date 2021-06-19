// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package access

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Data["authTypes"] = serverconfigs.FindAllHTTPAuthTypes()
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name string
	Type string

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
	params.Must.
		Field("name", params.Name).
		Require("请输入名称").
		Field("type", params.Type).
		Require("请输入认证类型")

	var ref = &serverconfigs.HTTPAuthPolicyRef{IsOn: true}
	var paramsJSON []byte

	switch params.Type {
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
		this.Fail("不支持的认证类型'" + params.Type + "'")
	}

	var paramsMap map[string]interface{}
	err := json.Unmarshal(paramsJSON, &paramsMap)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	createResp, err := this.RPC().HTTPAuthPolicyRPC().CreateHTTPAuthPolicy(this.AdminContext(), &pb.CreateHTTPAuthPolicyRequest{
		Name:       params.Name,
		Type:       params.Type,
		ParamsJSON: paramsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	defer this.CreateLogInfo("创建HTTP认证 %d", createResp.HttpAuthPolicyId)
	ref.AuthPolicyId = createResp.HttpAuthPolicyId
	ref.AuthPolicy = &serverconfigs.HTTPAuthPolicy{
		Id:     createResp.HttpAuthPolicyId,
		Name:   params.Name,
		IsOn:   true,
		Type:   params.Type,
		Params: paramsMap,
	}

	this.Data["policyRef"] = ref
	this.Success()
}
