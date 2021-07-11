// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package boards

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type UserAction struct {
	actionutils.ParentAction
}

func (this *UserAction) Init() {
	this.Nav("", "", "user")
}

func (this *UserAction) RunGet(params struct{}) {
	if !teaconst.IsPlus {
		this.RedirectURL("/dashboard")
		return
	}

	resp, err := this.RPC().UserRPC().ComposeUserGlobalBoard(this.AdminContext(), &pb.ComposeUserGlobalBoardRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["board"] = maps.Map{
		"totalUsers":            resp.TotalUsers,
		"countTodayUsers":       resp.CountTodayUsers,
		"countWeeklyUsers":      resp.CountWeeklyUsers,
		"countUserNodes":        resp.CountUserNodes,
		"countOfflineUserNodes": resp.CountOfflineUserNodes,
	}

	{
		statMaps := []maps.Map{}
		for _, stat := range resp.DailyStats {
			statMaps = append(statMaps, maps.Map{
				"day":   stat.Day,
				"count": stat.Count,
			})
		}
		this.Data["dailyStats"] = statMaps
	}

	// CPU
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.CpuNodeValues {
			statMaps = append(statMaps, maps.Map{
				"time":  timeutil.FormatTime("H:i", stat.CreatedAt),
				"value": types.Float32(string(stat.ValueJSON)),
			})
		}
		this.Data["cpuValues"] = statMaps
	}

	// Memory
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.MemoryNodeValues {
			statMaps = append(statMaps, maps.Map{
				"time":  timeutil.FormatTime("H:i", stat.CreatedAt),
				"value": types.Float32(string(stat.ValueJSON)),
			})
		}
		this.Data["memoryValues"] = statMaps
	}

	// Load
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.LoadNodeValues {
			statMaps = append(statMaps, maps.Map{
				"time":  timeutil.FormatTime("H:i", stat.CreatedAt),
				"value": types.Float32(string(stat.ValueJSON)),
			})
		}
		this.Data["loadValues"] = statMaps
	}

	// 流量排行
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.TopTrafficStats {
			statMaps = append(statMaps, maps.Map{
				"userId":        stat.UserId,
				"userName":      stat.UserName,
				"countRequests": stat.CountRequests,
				"bytes":         stat.Bytes,
			})
		}
		this.Data["topTrafficStats"] = statMaps
	}

	this.Show()
}
