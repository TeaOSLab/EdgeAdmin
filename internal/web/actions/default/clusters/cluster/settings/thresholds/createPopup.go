// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package thresholds

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	ClusterId int64
	NodeId    int64
}) {
	this.Data["clusterId"] = params.ClusterId
	this.Data["nodeId"] = params.NodeId
	this.Data["items"] = nodeconfigs.FindAllNodeValueItemDefinitions()
	this.Data["operators"] = nodeconfigs.FindAllNodeValueOperatorDefinitions()

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	ClusterId      int64
	NodeId         int64
	Item           string
	Param          string
	SumMethod      string
	Operator       string
	Value          string
	Duration       int32
	DurationUnit   string
	Message        string
	NotifyDuration int32

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	if params.ClusterId <= 0 && params.NodeId >= 0 {
		this.Fail("集群或者节点至少需要填写其中一个参数")
	}

	valueJSON, err := json.Marshal(params.Value)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	resp, err := this.RPC().NodeThresholdRPC().CreateNodeThreshold(this.AdminContext(), &pb.CreateNodeThresholdRequest{
		Role:           "node",
		NodeClusterId:  params.ClusterId,
		NodeId:         params.NodeId,
		Item:           params.Item,
		Param:          params.Param,
		Operator:       params.Operator,
		ValueJSON:      valueJSON,
		Message:        params.Message,
		Duration:       params.Duration,
		DurationUnit:   params.DurationUnit,
		SumMethod:      params.SumMethod,
		NotifyDuration: params.NotifyDuration,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	defer this.CreateLogInfo("创建节点阈值 %d", resp.NodeThresholdId)

	this.Success()
}
