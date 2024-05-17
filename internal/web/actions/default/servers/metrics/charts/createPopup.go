// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package charts

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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

func (this *CreatePopupAction) RunGet(params struct {
	ItemId int64
}) {
	this.Data["itemId"] = params.ItemId
	this.Data["types"] = serverconfigs.FindAllMetricChartTypes()

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	ItemId          int64
	Name            string
	Type            string
	WidthDiv        int32
	MaxItems        int32
	IgnoreEmptyKeys bool
	IgnoredKeys     []string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var chartId int64
	defer func() {
		this.CreateLogInfo(codes.MetricChart_LogCreateMetricChart, chartId)
	}()

	params.Must.
		Field("name", params.Name).
		Require("请输入图表名称").
		Field("type", params.Type).
		Require("请选择图表类型")

	createResp, err := this.RPC().MetricChartRPC().CreateMetricChart(this.AdminContext(), &pb.CreateMetricChartRequest{
		MetricItemId:    params.ItemId,
		Name:            params.Name,
		Type:            params.Type,
		WidthDiv:        params.WidthDiv,
		MaxItems:        params.MaxItems,
		ParamsJSON:      nil,
		IgnoreEmptyKeys: params.IgnoreEmptyKeys,
		IgnoredKeys:     params.IgnoredKeys,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	chartId = createResp.MetricChartId

	this.Success()
}
