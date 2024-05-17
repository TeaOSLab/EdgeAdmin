// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.
//go:build !plus

package access

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
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
	this.Data["authTypes"] = serverconfigs.FindAllHTTPAuthTypes(teaconst.Role)

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

	Exts        []string
	DomainsJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.HTTPAuthPolicy_LogUpdateHTTPAuthPolicy, params.PolicyId)

	policyResp, err := this.RPC().HTTPAuthPolicyRPC().FindEnabledHTTPAuthPolicy(this.AdminContext(), &pb.FindEnabledHTTPAuthPolicyRequest{HttpAuthPolicyId: params.PolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var policy = policyResp.HttpAuthPolicy
	if policy == nil {
		this.NotFound("httpAuthPolicy", params.PolicyId)
		return
	}
	policyType := policy.Type

	params.Must.
		Field("name", params.Name).
		Require("请输入名称")

	// 扩展名
	var exts = utils.NewStringsStream(params.Exts).
		Map(strings.TrimSpace, strings.ToLower).
		Filter(utils.FilterNotEmpty).
		Map(utils.MapAddPrefixFunc(".")).
		Unique().
		Result()

	// 域名
	var domains = []string{}
	if len(params.DomainsJSON) > 0 {
		var rawDomains = []string{}
		err := json.Unmarshal(params.DomainsJSON, &rawDomains)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// TODO 如果用户填写了一个网址，应该分析域名并填入

		domains = utils.NewStringsStream(rawDomains).
			Map(strings.TrimSpace, strings.ToLower).
			Filter(utils.FilterNotEmpty).
			Unique().
			Result()
	}

	var ref = &serverconfigs.HTTPAuthPolicyRef{IsOn: true}

	var method serverconfigs.HTTPAuthMethodInterface

	switch policyType {
	case serverconfigs.HTTPAuthTypeBasicAuth:
		var users = []*serverconfigs.HTTPAuthBasicMethodUser{}
		err := json.Unmarshal(params.HttpAuthBasicAuthUsersJSON, &users)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(users) == 0 {
			this.Fail("请添加至少一个用户")
		}
		method = &serverconfigs.HTTPAuthBasicMethod{
			Users:   users,
			Realm:   params.BasicAuthRealm,
			Charset: params.BasicAuthCharset,
		}
	case serverconfigs.HTTPAuthTypeSubRequest:
		params.Must.Field("subRequestURL", params.SubRequestURL).
			Require("请输入子请求URL")
		if params.SubRequestFollowRequest {
			params.SubRequestMethod = ""
		}
		method = &serverconfigs.HTTPAuthSubRequestMethod{
			URL:    params.SubRequestURL,
			Method: params.SubRequestMethod,
		}
	default:
		this.Fail("不支持的鉴权类型'" + policyType + "'")
	}

	if method == nil {
		this.Fail("无法找到对应的鉴权方式")
	}
	method.SetExts(exts)
	method.SetDomains(domains)

	paramsJSON, err := json.Marshal(method)
	if err != nil {
		this.ErrorPage(err)
		return
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
