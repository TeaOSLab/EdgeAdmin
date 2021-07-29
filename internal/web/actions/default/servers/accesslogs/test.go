// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipbox

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/accesslogs/policyutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type TestAction struct {
	actionutils.ParentAction
}

func (this *TestAction) Init() {
	this.Nav("", "", "test")
}

func (this *TestAction) RunGet(params struct {
	PolicyId int64
}) {
	err := policyutils.InitPolicy(this.Parent(), params.PolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}

func (this *TestAction) RunPost(params struct {
	PolicyId int64
	BodyJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo("测试向访问日志策略 %d 写入数据", params.PolicyId)

	var accessLog = &pb.HTTPAccessLog{}
	err := json.Unmarshal(params.BodyJSON, accessLog)
	if err != nil {
		this.Fail("发送内容不是有效的JSON：" + err.Error())
	}

	_, err = this.RPC().HTTPAccessLogPolicyRPC().WriteHTTPAccessLogPolicy(this.AdminContext(), &pb.WriteHTTPAccessLogPolicyRequest{
		HttpAccessLogPolicyId: params.PolicyId,
		HttpAccessLog:         accessLog,
	})
	if err != nil {
		this.Fail("发送失败：" + err.Error())
		return
	}

	this.Success()
}
