// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package metrics

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	ClusterId int64
	Category  string
}) {
	if len(params.Category) == 0 {
		params.Category = "http"
	}
	this.Data["category"] = params.Category
	this.Data["clusterId"] = params.ClusterId

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
		// 是否已添加
		existsResp, err := this.RPC().NodeClusterMetricItemRPC().ExistsNodeClusterMetricItem(this.AdminContext(), &pb.ExistsNodeClusterMetricItemRequest{
			NodeClusterId: params.ClusterId,
			MetricItemId:  item.Id,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var exists = existsResp.Exists

		itemMaps = append(itemMaps, maps.Map{
			"id":             item.Id,
			"name":           item.Name,
			"code":           item.Code,
			"isOn":           item.IsOn,
			"period":         item.Period,
			"periodUnit":     item.PeriodUnit,
			"periodUnitName": serverconfigs.FindMetricPeriodUnitName(item.PeriodUnit),
			"keys":           item.Keys,
			"value":          item.Value,
			"valueName":      serverconfigs.FindMetricValueName(item.Category, item.Value),
			"category":       item.Category,
			"isPublic":       item.IsPublic,
			"isChecked":      exists,
		})
	}
	this.Data["items"] = itemMaps

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	ClusterId int64
	ItemId    int64

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.MetricItem_LogAddMetricItemToCluster, params.ItemId, params.ClusterId)

	_, err := this.RPC().NodeClusterMetricItemRPC().EnableNodeClusterMetricItem(this.AdminContext(), &pb.EnableNodeClusterMetricItemRequest{
		NodeClusterId: params.ClusterId,
		MetricItemId:  params.ItemId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
