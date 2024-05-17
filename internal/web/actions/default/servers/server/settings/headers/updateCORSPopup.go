// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type UpdateCORSPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateCORSPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateCORSPopupAction) RunGet(params struct {
	HeaderPolicyId int64
}) {
	this.Data["headerPolicyId"] = params.HeaderPolicyId

	resp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{HttpHeaderPolicyId: params.HeaderPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var headerPolicyJSON = resp.HttpHeaderPolicyJSON
	var headerPolicy = &shared.HTTPHeaderPolicy{}
	if len(headerPolicyJSON) > 0 {
		err = json.Unmarshal(headerPolicyJSON, headerPolicy)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["cors"] = headerPolicy.CORS

	this.Show()
}

func (this *UpdateCORSPopupAction) RunPost(params struct {
	HeaderPolicyId int64
	CorsJSON       []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var config = shared.NewHTTPCORSHeaderConfig()
	err := json.Unmarshal(params.CorsJSON, config)
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
		return
	}

	err = config.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
		return
	}

	_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicyCORS(this.AdminContext(), &pb.UpdateHTTPHeaderPolicyCORSRequest{
		HttpHeaderPolicyId: params.HeaderPolicyId,
		CorsJSON:           params.CorsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
