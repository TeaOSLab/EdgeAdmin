// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package metricutils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

// InitItem 初始化指标信息
func InitItem(parent *actionutils.ParentAction, itemId int64) (*pb.MetricItem, error) {
	client, err := rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := client.MetricItemRPC().FindEnabledMetricItem(parent.AdminContext(), &pb.FindEnabledMetricItemRequest{MetricItemId: itemId})
	if err != nil {
		return nil, err
	}
	var item = resp.MetricItem
	if item == nil {
		return nil, errors.New("metric item not found")
	}

	countChartsResp, err := client.MetricChartRPC().CountEnabledMetricCharts(parent.AdminContext(), &pb.CountEnabledMetricChartsRequest{MetricItemId: item.Id})
	if err != nil {
		return nil, err
	}
	var countCharts = countChartsResp.Count

	parent.Data["item"] = maps.Map{
		"id":             item.Id,
		"name":           item.Name,
		"code":           item.Code,
		"isOn":           item.IsOn,
		"keys":           item.Keys,
		"value":          item.Value,
		"valueName":      serverconfigs.FindMetricValueName(item.Category, item.Value),
		"period":         item.Period,
		"periodUnit":     item.PeriodUnit,
		"periodUnitName": serverconfigs.FindMetricPeriodUnitName(item.PeriodUnit),
		"category":       item.Category,
		"isPublic":       item.IsPublic,
		"countCharts":    countCharts,
		"expiresPeriod":  item.ExpiresPeriod,
	}
	return item, nil
}
