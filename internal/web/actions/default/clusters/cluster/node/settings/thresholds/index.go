// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package thresholds

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "update")
	this.SecondMenu("threshold")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
	NodeId    int64
}) {
	_, err := nodeutils.InitNodeInfo(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["nodeId"] = params.NodeId

	// 列出所有阈值
	thresholdsResp, err := this.RPC().NodeThresholdRPC().FindAllEnabledNodeThresholds(this.AdminContext(), &pb.FindAllEnabledNodeThresholdsRequest{
		Role:          "node",
		NodeClusterId: params.ClusterId,
		NodeId:        params.NodeId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	thresholdMaps := []maps.Map{}
	for _, threshold := range thresholdsResp.NodeThresholds {
		thresholdMaps = append(thresholdMaps, maps.Map{
			"id":               threshold.Id,
			"itemName":         nodeconfigs.FindNodeValueItemName(threshold.Item),
			"paramName":        nodeconfigs.FindNodeValueItemParamName(threshold.Item, threshold.Param),
			"operatorName":     nodeconfigs.FindNodeValueOperatorName(threshold.Operator),
			"value":            nodeconfigs.UnmarshalNodeValue(threshold.ValueJSON),
			"sumMethodName":    nodeconfigs.FindNodeValueSumMethodName(threshold.SumMethod),
			"duration":         threshold.Duration,
			"durationUnitName": nodeconfigs.FindNodeValueDurationUnitName(threshold.DurationUnit),
			"isOn":             threshold.IsOn,
		})
	}
	this.Data["thresholds"] = thresholdMaps

	this.Show()
}
