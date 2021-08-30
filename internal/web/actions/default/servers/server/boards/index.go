package boards

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	"strconv"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "board", "")
	this.SecondMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	if !teaconst.IsPlus {
		this.RedirectURL("/servers/server/stat?serverId=" + strconv.FormatInt(params.ServerId, 10))
		return
	}

	serverResp, err := this.RPC().ServerRPC().FindEnabledServer(this.AdminContext(), &pb.FindEnabledServerRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var server = serverResp.Server
	if server == nil {
		this.NotFound("server", params.ServerId)
		return
	}
	this.Data["server"] = maps.Map{
		"id":   server.Id,
		"name": server.Name,
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId int64
}) {
	resp, err := this.RPC().ServerStatBoardRPC().ComposeServerStatBoard(this.AdminContext(), &pb.ComposeServerStatBoardRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
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

	// 节点排行
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.TopNodeStats {
			statMaps = append(statMaps, maps.Map{
				"nodeId":              stat.NodeId,
				"nodeName":            stat.NodeName,
				"countRequests":       stat.CountRequests,
				"bytes":               stat.Bytes,
				"countAttackRequests": stat.CountAttackRequests,
				"attackBytes":         stat.AttackBytes,
			})
		}
		this.Data["topNodeStats"] = statMaps
	}

	// 域名排行
	{
		var statMaps = []maps.Map{}
		for _, stat := range resp.TopDomainStats {
			statMaps = append(statMaps, maps.Map{
				"serverId":            stat.ServerId,
				"domain":              stat.Domain,
				"countRequests":       stat.CountRequests,
				"bytes":               stat.Bytes,
				"countAttackRequests": stat.CountAttackRequests,
				"attackBytes":         stat.AttackBytes,
			})
		}
		this.Data["topDomainStats"] = statMaps
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
	this.Success()
}
