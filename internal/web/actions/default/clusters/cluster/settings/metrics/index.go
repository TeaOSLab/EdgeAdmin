// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package metrics

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "setting")
	this.SecondMenu("metric")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
	Category  string
}) {
	if len(params.Category) == 0 {
		params.Category = "http"
	}
	this.Data["category"] = params.Category

	itemsResp, err := this.RPC().NodeClusterMetricItemRPC().FindAllNodeClusterMetricItems(this.AdminContext(), &pb.FindAllNodeClusterMetricItemsRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var itemMaps = []maps.Map{}
	for _, item := range itemsResp.MetricItems {
		itemMaps = append(itemMaps, maps.Map{
			"id":         item.Id,
			"name":       item.Name,
			"isOn":       item.IsOn,
			"period":     item.Period,
			"periodUnit": item.PeriodUnit,
			"keys":       item.Keys,
			"value":      item.Value,
			"category":   item.Category,
		})
	}
	this.Data["items"] = itemMaps

	this.Show()
}
