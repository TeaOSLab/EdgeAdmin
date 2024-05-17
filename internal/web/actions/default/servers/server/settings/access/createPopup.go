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
	"strings"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Data["authTypes"] = serverconfigs.FindAllHTTPAuthTypes(teaconst.Role)
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

	Exts        []string
	DomainsJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入名称").
		Field("type", params.Type).
		Require("请输入鉴权类型")

	var ref = &serverconfigs.HTTPAuthPolicyRef{IsOn: true}

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

	var method serverconfigs.HTTPAuthMethodInterface

	switch params.Type {
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
		this.Fail("不支持的鉴权类型'" + params.Type + "'")
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

	createResp, err := this.RPC().HTTPAuthPolicyRPC().CreateHTTPAuthPolicy(this.AdminContext(), &pb.CreateHTTPAuthPolicyRequest{
		Name:       params.Name,
		Type:       params.Type,
		ParamsJSON: paramsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	defer this.CreateLogInfo(codes.HTTPAuthPolicy_LogCreateHTTPAuthPolicy, createResp.HttpAuthPolicyId)
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
