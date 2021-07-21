// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package boards

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"strconv"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "board")
	this.SecondMenu("nodes")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
	NodeId    int64
}) {
	err := nodeutils.InitNodeInfo(this, params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if !teaconst.IsPlus {
		this.RedirectURL("/clusters/cluster/node?clusterId=" + strconv.FormatInt(params.ClusterId, 10) + "&nodeId=" + strconv.FormatInt(params.NodeId, 10))
		return
	}

	resp, err := this.RPC().ServerStatBoardRPC().ComposeServerStatNodeBoard(this.AdminContext(), &pb.ComposeServerStatNodeBoardRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["board"] = maps.Map{
		"isActive":            resp.IsActive,
		"trafficInBytes":      resp.TrafficInBytes,
		"trafficOutBytes":     resp.TrafficOutBytes,
		"countConnections":    resp.CountConnections,
		"countRequests":       resp.CountRequests,
		"countAttackRequests": resp.CountAttackRequests,
		"cpuUsage":            resp.CpuUsage,
		"memoryUsage":         resp.MemoryUsage,
		"memoryTotalSize":     resp.MemoryTotalSize,
		"load":                resp.Load,
		"cacheDiskSize":       resp.CacheDiskSize,
		"cacheMemorySize":     resp.CacheMemorySize,
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
				"countAttackRequests": stat.CountAttackRequests,
				"attackBytes":         stat.AttackBytes,
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
				"countAttackRequests": stat.CountAttackRequests,
				"attackBytes":         stat.AttackBytes,
				"day":                 stat.Day[4:6] + "月" + stat.Day[6:] + "日",
			})
		}
		this.Data["dailyStats"] = statMaps
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

	// CacheDirs
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.CacheDirsValues {
			var m = maps.Map{}
			err = json.Unmarshal(stat.ValueJSON, &m)
			if err != nil {
				continue
			}
			statMaps = append(statMaps, maps.Map{
				"time":  timeutil.FormatTime("H:i", stat.CreatedAt),
				"value": m,
			})
		}
		this.Data["cacheDirValues"] = statMaps
	}

	// 指标
	{
		var chartMaps = []maps.Map{}
		for _, chart := range resp.MetricDataCharts {
			var statMaps = []maps.Map{}
			for _, stat := range chart.MetricStats {
				statMaps = append(statMaps, maps.Map{
					"keys":  stat.Keys,
					"time":  stat.Time,
					"value": stat.Value,
					"count": stat.SumCount,
					"total": stat.SumTotal,
				})
			}
			chartMaps = append(chartMaps, maps.Map{
				"chart": maps.Map{
					"id":       chart.MetricChart.Id,
					"name":     chart.MetricChart.Name,
					"widthDiv": chart.MetricChart.WidthDiv,
					"isOn":     chart.MetricChart.IsOn,
					"maxItems": chart.MetricChart.MaxItems,
					"type":     chart.MetricChart.Type,
				},
				"item": maps.Map{
					"id":            chart.MetricChart.MetricItem.Id,
					"name":          chart.MetricChart.MetricItem.Name,
					"period":        chart.MetricChart.MetricItem.Period,
					"periodUnit":    chart.MetricChart.MetricItem.PeriodUnit,
					"valueType":     serverconfigs.FindMetricValueType(chart.MetricChart.MetricItem.Category, chart.MetricChart.MetricItem.Value),
					"valueTypeName": serverconfigs.FindMetricValueName(chart.MetricChart.MetricItem.Category, chart.MetricChart.MetricItem.Value),
					"keys":          chart.MetricChart.MetricItem.Keys,
				},
				"stats": statMaps,
			})
		}
		this.Data["metricCharts"] = chartMaps
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId int64
	NodeId    int64
}) {
	resp, err := this.RPC().ServerStatBoardRPC().ComposeServerStatNodeBoard(this.AdminContext(), &pb.ComposeServerStatNodeBoardRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["board"] = maps.Map{
		"isActive":            resp.IsActive,
		"trafficInBytes":      resp.TrafficInBytes,
		"trafficOutBytes":     resp.TrafficOutBytes,
		"countConnections":    resp.CountConnections,
		"countRequests":       resp.CountRequests,
		"countAttackRequests": resp.CountAttackRequests,
		"cpuUsage":            resp.CpuUsage,
		"memoryUsage":         resp.MemoryUsage,
		"memoryTotalSize":     resp.MemoryTotalSize,
		"load":                resp.Load,
		"cacheDiskSize":       resp.CacheDiskSize,
		"cacheMemorySize":     resp.CacheMemorySize,
	}

	this.Success()
}
