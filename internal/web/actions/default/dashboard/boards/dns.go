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

type DnsAction struct {
	actionutils.ParentAction
}

func (this *DnsAction) Init() {
	this.Nav("", "", "dns")
}

func (this *DnsAction) RunGet(params struct{}) {
	if !teaconst.IsPlus {
		this.RedirectURL("/dashboard")
		return
	}

	resp, err := this.RPC().NSRPC().ComposeNSBoard(this.AdminContext(), &pb.ComposeNSBoardRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["board"] = maps.Map{
		"countDomains":      resp.CountNSDomains,
		"countRecords":      resp.CountNSRecords,
		"countClusters":     resp.CountNSClusters,
		"countNodes":        resp.CountNSNodes,
		"countOfflineNodes": resp.CountOfflineNSNodes,
	}

	// 流量排行
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.HourlyTrafficStats {
			statMaps = append(statMaps, maps.Map{
				"day":           stat.Hour[4:6] + "月" + stat.Hour[6:8] + "日",
				"hour":          stat.Hour[8:],
				"countRequests": stat.CountRequests,
				"bytes":         stat.Bytes,
			})
		}
		this.Data["hourlyStats"] = statMaps
	}

	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.DailyTrafficStats {
			statMaps = append(statMaps, maps.Map{
				"day":           stat.Day[4:6] + "月" + stat.Day[6:] + "日",
				"countRequests": stat.CountRequests,
				"bytes":         stat.Bytes,
			})
		}
		this.Data["dailyStats"] = statMaps
	}

	// 域名排行
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.TopNSDomainStats {
			statMaps = append(statMaps, maps.Map{
				"domainId":      stat.NsDomainId,
				"domainName":    stat.NsDomainName,
				"countRequests": stat.CountRequests,
				"bytes":         stat.Bytes,
			})
		}
		this.Data["topDomainStats"] = statMaps
	}

	// 节点排行
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.TopNSNodeStats {
			statMaps = append(statMaps, maps.Map{
				"clusterId":     stat.NsClusterId,
				"nodeId":        stat.NsNodeId,
				"nodeName":      stat.NsNodeName,
				"countRequests": stat.CountRequests,
				"bytes":         stat.Bytes,
			})
		}
		this.Data["topNodeStats"] = statMaps
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
