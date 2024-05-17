// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package stat

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"sort"
)

type HourlyRequestsAction struct {
	actionutils.ParentAction
}

func (this *HourlyRequestsAction) Init() {
	this.Nav("", "stat", "hourly")
	this.SecondMenu("index")
}

func (this *HourlyRequestsAction) RunGet(params struct {
	ServerId int64
}) {
	this.Data["serverId"] = params.ServerId

	resp, err := this.RPC().ServerDailyStatRPC().FindLatestServerHourlyStats(this.AdminContext(), &pb.FindLatestServerHourlyStatsRequest{
		ServerId: params.ServerId,
		Hours:    24,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	sort.Slice(resp.Stats, func(i, j int) bool {
		stat1 := resp.Stats[i]
		stat2 := resp.Stats[j]
		return stat1.Hour < stat2.Hour
	})
	statMaps := []maps.Map{}
	for _, stat := range resp.Stats {
		statMaps = append(statMaps, maps.Map{
			"day":                 stat.Hour[:4] + "-" + stat.Hour[4:6] + "-" + stat.Hour[6:8],
			"hour":                stat.Hour[8:],
			"bytes":               stat.Bytes,
			"cachedBytes":         stat.CachedBytes,
			"countRequests":       stat.CountRequests,
			"countCachedRequests": stat.CountCachedRequests,
		})
	}
	this.Data["hourlyStats"] = statMaps

	this.Show()
}
