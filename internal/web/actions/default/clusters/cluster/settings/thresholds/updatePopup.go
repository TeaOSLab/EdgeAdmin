// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package thresholds

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
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
	ThresholdId int64
}) {
	// 通用参数
	this.Data["items"] = nodeconfigs.FindAllNodeValueItemDefinitions()
	this.Data["operators"] = nodeconfigs.FindAllNodeValueOperatorDefinitions()

	// 阈值详情
	thresholdResp, err := this.RPC().NodeThresholdRPC().FindEnabledNodeThreshold(this.AdminContext(), &pb.FindEnabledNodeThresholdRequest{NodeThresholdId: params.ThresholdId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	threshold := thresholdResp.NodeThreshold
	if threshold == nil {
		this.NotFound("nodeThreshold", params.ThresholdId)
		return
	}

	valueInterface := new(interface{})
	err = json.Unmarshal(threshold.ValueJSON, valueInterface)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["threshold"] = maps.Map{
		"id":             threshold.Id,
		"item":           threshold.Item,
		"param":          threshold.Param,
		"message":        threshold.Message,
		"notifyDuration": threshold.NotifyDuration,
		"value":          nodeconfigs.UnmarshalNodeValue(threshold.ValueJSON),
		"operator":       threshold.Operator,
		"duration":       threshold.Duration,
		"durationUnit":   threshold.DurationUnit,
		"isOn":           threshold.IsOn,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	ThresholdId    int64
	Item           string
	Param          string
	SumMethod      string
	Operator       string
	Value          string
	Duration       int32
	DurationUnit   string
	Message        string
	NotifyDuration int32
	IsOn           bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改节点阈值 %d", params.ThresholdId)

	valueJSON, err := json.Marshal(params.Value)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().NodeThresholdRPC().UpdateNodeThreshold(this.AdminContext(), &pb.UpdateNodeThresholdRequest{
		NodeThresholdId: params.ThresholdId,
		Item:            params.Item,
		Param:           params.Param,
		Operator:        params.Operator,
		ValueJSON:       valueJSON,
		Message:         params.Message,
		NotifyDuration:  params.NotifyDuration,
		Duration:        params.Duration,
		DurationUnit:    params.DurationUnit,
		SumMethod:       params.SumMethod,
		IsOn:            params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
