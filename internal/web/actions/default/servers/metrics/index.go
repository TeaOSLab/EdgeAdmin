// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package metrics

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct {
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
			"id":             item.Id,
			"name":           item.Name,
			"code":           item.Code,
			"isOn":           item.IsOn,
			"period":         item.Period,
			"periodUnit":     item.PeriodUnit,
			"periodUnitName": serverconfigs.FindMetricPeriodUnitName(item.PeriodUnit),
			"expiresPeriod":  item.ExpiresPeriod,
			"keys":           item.Keys,
			"value":          item.Value,
			"valueName":      serverconfigs.FindMetricValueName(item.Category, item.Value),
			"category":       item.Category,
			"isPublic":       item.IsPublic,
		})
	}
	this.Data["items"] = itemMaps

	this.Show()
}
