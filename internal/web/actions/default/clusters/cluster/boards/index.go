// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package boards

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "board", "")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	resp, err := this.RPC().ServerStatBoardRPC().ComposeServerStatNodeClusterBoard(this.AdminContext(), &pb.ComposeServerStatNodeClusterBoardRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["board"] = maps.Map{
		"countUsers":         resp.CountUsers,
		"countActiveNodes":   resp.CountActiveNodes,
		"countInactiveNodes": resp.CountInactiveNodes,
		"countServers":       resp.CountServers,
	}

	// 24小时流量趋势
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.HourlyTrafficStats {
			statMaps = append(statMaps, maps.Map{
				"bytes":               stat.Bytes,
				"cachedBytes":         stat.CachedBytes,
				"countRequests":       stat.CountRequests,
				"countCachedRequests": stat.CountCachedRequests,
				"day":                 stat.Hour[4:6] + "月" + stat.Hour[6:8] + "日",
				"hour":                stat.Hour[8:],
			})
		}
		this.Data["hourlyStats"] = statMaps
	}

	// 15天流量趋势
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.DailyTrafficStats {
			statMaps = append(statMaps, maps.Map{
				"bytes":               stat.Bytes,
				"cachedBytes":         stat.CachedBytes,
				"countRequests":       stat.CountRequests,
				"countCachedRequests": stat.CountCachedRequests,
				"day":                 stat.Day[4:6] + "月" + stat.Day[6:] + "日",
			})
		}
		this.Data["dailyStats"] = statMaps
	}

	// 节点排行
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.TopNodeStats {
			statMaps = append(statMaps, maps.Map{
				"nodeId":        stat.NodeId,
				"nodeName":      stat.NodeName,
				"countRequests": stat.CountRequests,
				"bytes":         stat.Bytes,
			})
		}
		this.Data["topNodeStats"] = statMaps
	}

	// 域名排行
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.TopDomainStats {
			statMaps = append(statMaps, maps.Map{
				"serverId":      stat.ServerId,
				"domain":        stat.Domain,
				"countRequests": stat.CountRequests,
				"bytes":         stat.Bytes,
			})
		}
		this.Data["topDomainStats"] = statMaps
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

	this.Show()
}
