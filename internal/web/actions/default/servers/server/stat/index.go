// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package stat

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"sort"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "stat", "minutely")
	this.SecondMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	this.Data["serverId"] = params.ServerId

	resp, err := this.RPC().ServerDailyStatRPC().FindLatestServerMinutelyStats(this.AdminContext(), &pb.FindLatestServerMinutelyStatsRequest{
		ServerId: params.ServerId,
		Minutes:  120,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	sort.Slice(resp.Stats, func(i, j int) bool {
		stat1 := resp.Stats[i]
		stat2 := resp.Stats[j]
		return stat1.Minute < stat2.Minute
	})
	statMaps := []maps.Map{}
	for _, stat := range resp.Stats {
		statMaps = append(statMaps, maps.Map{
			"day":                 stat.Minute[:4] + "-" + stat.Minute[4:6] + "-" + stat.Minute[6:8],
			"minute":              stat.Minute[8:10] + ":" + stat.Minute[10:12],
			"bytes":               stat.Bytes,
			"cachedBytes":         stat.CachedBytes,
			"countRequests":       stat.CountRequests,
			"countCachedRequests": stat.CountCachedRequests,
		})
	}
	this.Data["minutelyStats"] = statMaps

	this.Show()
}
