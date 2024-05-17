// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package api

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"sort"
	"strings"
)

type MethodStatsAction struct {
	actionutils.ParentAction
}

func (this *MethodStatsAction) Init() {
	this.Nav("", "", "")
}

func (this *MethodStatsAction) RunGet(params struct {
	Order  string
	Method string
	Tag    string
}) {
	this.Data["order"] = params.Order
	this.Data["method"] = params.Method
	this.Data["tag"] = params.Tag

	statsResp, err := this.RPC().APIMethodStatRPC().FindAPIMethodStatsWithDay(this.AdminContext(), &pb.FindAPIMethodStatsWithDayRequest{Day: timeutil.Format("Ymd")})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var pbStats = statsResp.ApiMethodStats

	switch params.Order {
	case "method":
		sort.Slice(pbStats, func(i, j int) bool {
			return pbStats[i].Method < pbStats[j].Method
		})
	case "costMs.desc":
		sort.Slice(pbStats, func(i, j int) bool {
			return pbStats[i].CostMs > pbStats[j].CostMs
		})
	case "peekMs.desc":
		sort.Slice(pbStats, func(i, j int) bool {
			return pbStats[i].PeekMs > pbStats[j].PeekMs
		})
	case "calls.desc":
		sort.Slice(pbStats, func(i, j int) bool {
			return pbStats[i].CountCalls > pbStats[j].CountCalls
		})
	}

	var statMaps = []maps.Map{}
	for _, stat := range pbStats {
		if len(params.Method) > 0 && !strings.Contains(strings.ToLower(stat.Method), strings.ToLower(params.Method)) {
			continue
		}
		if len(params.Tag) > 0 && !strings.Contains(strings.ToLower(stat.Tag), strings.ToLower(params.Tag)) {
			continue
		}

		statMaps = append(statMaps, maps.Map{
			"id":         stat.Id,
			"method":     stat.Method,
			"tag":        stat.Tag,
			"costMs":     stat.CostMs,
			"peekMs":     stat.PeekMs,
			"countCalls": stat.CountCalls,
		})
	}
	this.Data["stats"] = statMaps

	this.Show()
}
