// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package dashboard

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dashboard/dashboardutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	"regexp"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	// 通知菜单数字Badge更新
	helpers.NotifyIPItemsCountChanges()
	helpers.NotifyNodeLogsCountChange()

	if this.checkPlus() {
		this.RedirectURL("/dashboard/boards")
		return
	}

	// 取得用户的权限
	module, ok := configloaders.FindFirstAdminModule(this.AdminId())
	if ok {
		if module != "dashboard" {
			for _, m := range configloaders.AllModuleMaps(this.LangCode()) {
				if m.GetString("code") == module {
					this.RedirectURL(m.GetString("url"))
					return
				}
			}
		}
	}

	// 版本更新
	this.Data["currentVersionCode"] = teaconst.Version
	this.Data["newVersionCode"] = teaconst.NewVersionCode
	this.Data["newVersionDownloadURL"] = teaconst.NewVersionDownloadURL

	this.Show()
}

func (this *IndexAction) RunPost(params struct{}) {
	// 读取看板数据
	resp, err := this.RPC().AdminRPC().ComposeAdminDashboard(this.AdminContext(), &pb.ComposeAdminDashboardRequest{
		ApiVersion: teaconst.APINodeVersion,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 检查当前服务器空间
	var diskUsageWarning = ""
	diskPath, diskUsage, diskUsagePercent, shouldWarning := dashboardutils.CheckDiskPartitions(90)
	if shouldWarning {
		diskUsageWarning = codes.AdminDashboard_DiskUsageWarning.For(this.LangCode(), diskPath, diskUsage/(1<<30), diskUsagePercent, 100-diskUsagePercent)
	}

	this.Data["dashboard"] = maps.Map{
		"defaultClusterId": resp.DefaultNodeClusterId,

		"countServers":         resp.CountServers,
		"countNodeClusters":    resp.CountNodeClusters,
		"countNodes":           resp.CountNodes,
		"countOfflineNodes":    resp.CountOfflineNodes,
		"countUsers":           resp.CountUsers,
		"countAPINodes":        resp.CountAPINodes,
		"countOfflineAPINodes": resp.CountOfflineAPINodes,
		"countDBNodes":         resp.CountDBNodes,

		"canGoServers":  configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeServer),
		"canGoNodes":    configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeNode),
		"canGoSettings": configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeSetting),
		"canGoUsers":    configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeUser),

		"diskUsageWarning": diskUsageWarning,
	}

	// 今日流量和独立IP数
	var todayTrafficBytes int64
	var todayCountIPs int64
	if len(resp.DailyTrafficStats) > 0 {
		var lastDailyTrafficStat = resp.DailyTrafficStats[len(resp.DailyTrafficStats)-1]
		todayTrafficBytes = lastDailyTrafficStat.Bytes
		todayCountIPs = lastDailyTrafficStat.CountIPs
	}
	var todayTrafficString = numberutils.FormatBytes(todayTrafficBytes)
	var result = regexp.MustCompile(`^(?U)(.+)([a-zA-Z]+)$`).FindStringSubmatch(todayTrafficString)
	if len(result) > 2 {
		this.Data["todayTraffic"] = result[1]
		this.Data["todayTrafficUnit"] = result[2]
	} else {
		this.Data["todayTraffic"] = todayTrafficString
		this.Data["todayTrafficUnit"] = ""
	}

	this.Data["todayCountIPs"] = todayCountIPs

	// 24小时流量趋势
	{
		statMaps := []maps.Map{}
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
		this.Data["hourlyTrafficStats"] = statMaps
	}

	// 15天流量趋势
	{
		statMaps := []maps.Map{}
		for _, stat := range resp.DailyTrafficStats {
			statMaps = append(statMaps, maps.Map{
				"bytes":               stat.Bytes,
				"cachedBytes":         stat.CachedBytes,
				"countRequests":       stat.CountRequests,
				"countCachedRequests": stat.CountCachedRequests,
				"countAttackRequests": stat.CountAttackRequests,
				"attackBytes":         stat.AttackBytes,
				"day":                 stat.Day[4:6] + "月" + stat.Day[6:] + "日",
				"countIPs":            stat.CountIPs,
			})
		}
		this.Data["dailyTrafficStats"] = statMaps
	}

	// 版本升级
	if resp.NodeUpgradeInfo != nil {
		this.Data["nodeUpgradeInfo"] = maps.Map{
			"count":   resp.NodeUpgradeInfo.CountNodes,
			"version": resp.NodeUpgradeInfo.NewVersion,
		}
	} else {
		this.Data["nodeUpgradeInfo"] = maps.Map{
			"count":   0,
			"version": "",
		}
	}
	if resp.ApiNodeUpgradeInfo != nil {
		this.Data["apiNodeUpgradeInfo"] = maps.Map{
			"count":   resp.ApiNodeUpgradeInfo.CountNodes,
			"version": resp.ApiNodeUpgradeInfo.NewVersion,
		}
	} else {
		this.Data["apiNodeUpgradeInfo"] = maps.Map{
			"count":   0,
			"version": "",
		}
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

	// 当前API节点版本
	{
		exePath, runtimeVersion, fileVersion, ok := dashboardutils.CheckLocalAPINode(this.RPC(), this.AdminContext())
		if ok {
			this.Data["localLowerVersionAPINode"] = maps.Map{
				"exePath":        exePath,
				"runtimeVersion": runtimeVersion,
				"fileVersion":    fileVersion,
				"isRestarting":   false,
			}
		}
	}

	// 弱密码提示
	countWeakAdminsResp, err := this.RPC().AdminRPC().CountAllEnabledAdmins(this.AdminContext(), &pb.CountAllEnabledAdminsRequest{HasWeakPassword: true})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countWeakAdmins"] = countWeakAdminsResp.Count

	this.Success()
}
