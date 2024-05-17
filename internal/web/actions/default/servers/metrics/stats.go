// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package metrics

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/metrics/metricutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

type StatsAction struct {
	actionutils.ParentAction
}

func (this *StatsAction) Init() {
	this.Nav("", "", "stat")
}

func (this *StatsAction) RunGet(params struct {
	ItemId int64
}) {
	item, err := metricutils.InitItem(this.Parent(), params.ItemId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	countResp, err := this.RPC().MetricStatRPC().CountMetricStats(this.AdminContext(), &pb.CountMetricStatsRequest{
		MetricItemId: params.ItemId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count

	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	statsResp, err := this.RPC().MetricStatRPC().ListMetricStats(this.AdminContext(), &pb.ListMetricStatsRequest{
		MetricItemId: params.ItemId,
		Offset:       page.Offset,
		Size:         page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var statMaps = []maps.Map{}
	for _, stat := range statsResp.MetricStats {
		// 占比
		var ratio float32
		if stat.SumTotal > 0 {
			ratio = stat.Value * 100 / stat.SumTotal
		}

		statMaps = append(statMaps, maps.Map{
			"id":      stat.Id,
			"time":    serverconfigs.HumanMetricTime(item.PeriodUnit, stat.Time),
			"keys":    stat.Keys,
			"value":   stat.Value,
			"cluster": maps.Map{"id": stat.NodeCluster.Id, "name": stat.NodeCluster.Name},
			"node":    maps.Map{"id": stat.Node.Id, "name": stat.Node.Name},
			"server":  maps.Map{"id": stat.Server.Id, "name": stat.Server.Name},
			"ratio":   fmt.Sprintf("%.2f", ratio),
		})
	}
	this.Data["stats"] = statMaps

	this.Show()
}
