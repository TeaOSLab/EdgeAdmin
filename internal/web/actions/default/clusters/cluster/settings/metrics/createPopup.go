// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package metrics

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	Category string
}) {
	if len(params.Category) == 0 {
		params.Category = "http"
	}
	this.Data["category"] = params.Category

	countResp, err := this.RPC().MetricItemRPC().CountAllEnabledMetricItems(this.AdminContext(), &pb.CountAllEnabledMetricItemsRequest{Category: params.Category})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	itemsResp, err := this.RPC().MetricItemRPC().ListEnabledMetricItems(this.AdminContext(), &pb.ListEnabledMetricItemsRequest{
		Category: params.Category,
		Offset:   page.Offset,
		Size:     page.Size,
	})
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
