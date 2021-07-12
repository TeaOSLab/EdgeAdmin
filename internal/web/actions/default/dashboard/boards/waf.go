// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package boards

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type WafAction struct {
	actionutils.ParentAction
}

func (this *WafAction) Init() {
	this.Nav("", "", "waf")
}

func (this *WafAction) RunGet(params struct{}) {
	if !teaconst.IsPlus {
		this.RedirectURL("/dashboard")
		return
	}

	resp, err := this.RPC().FirewallRPC().ComposeFirewallGlobalBoard(this.AdminContext(), &pb.ComposeFirewallGlobalBoardRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["board"] = maps.Map{
		"countDailyLogs":    resp.CountDailyLogs,
		"countDailyBlocks":  resp.CountDailyBlocks,
		"countDailyCaptcha": resp.CountDailyCaptcha,
		"countWeeklyBlocks": resp.CountWeeklyBlocks,
	}

	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.HourlyStats {
			statMaps = append(statMaps, maps.Map{
				"hour":         stat.Hour,
				"countLogs":    stat.CountLogs,
				"countCaptcha": stat.CountCaptcha,
				"countBlocks":  stat.CountBlocks,
			})
		}
		this.Data["hourlyStats"] = statMaps
	}

	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.DailyStats {
			statMaps = append(statMaps, maps.Map{
				"day":          stat.Day,
				"countLogs":    stat.CountLogs,
				"countCaptcha": stat.CountCaptcha,
				"countBlocks":  stat.CountBlocks,
			})
		}
		this.Data["dailyStats"] = statMaps
	}

	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.HttpFirewallRuleGroups {
			statMaps = append(statMaps, maps.Map{
				"name":  stat.HttpFirewallRuleGroup.Name,
				"count": stat.Count,
			})
		}
		this.Data["groupStats"] = statMaps
	}

	this.Show()
}
