// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package chartutils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

// InitChart 初始化指标图表信息
func InitChart(parent *actionutils.ParentAction, chartId int64) (*pb.MetricChart, error) {
	client, err := rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := client.MetricChartRPC().FindEnabledMetricChart(parent.AdminContext(), &pb.FindEnabledMetricChartRequest{MetricChartId: chartId})
	if err != nil {
		return nil, err
	}
	var chart = resp.MetricChart
	if chart == nil {
		return nil, errors.New("metric chart not found")
	}
	if len(chart.IgnoredKeys) == 0 {
		chart.IgnoredKeys = []string{}
	}
	parent.Data["chart"] = maps.Map{
		"id":              chart.Id,
		"name":            chart.Name,
		"isOn":            chart.IsOn,
		"widthDiv":        chart.WidthDiv,
		"maxItems":        chart.MaxItems,
		"type":            chart.Type,
		"typeName":        serverconfigs.FindMetricChartTypeName(chart.Type),
		"ignoreEmptyKeys": chart.IgnoreEmptyKeys,
		"ignoredKeys":     chart.IgnoredKeys,
	}
	return chart, nil
}
