package dashboard

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"math"
	"regexp"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	// 取得用户的权限
	module, ok := configloaders.FindFirstAdminModule(this.AdminId())
	if ok {
		if module != "dashboard" {
			for _, m := range configloaders.AllModuleMaps() {
				if m.GetString("code") == module {
					this.RedirectURL(m.GetString("url"))
					return
				}
			}
		}
	}

	// 读取看板数据
	resp, err := this.RPC().AdminRPC().ComposeAdminDashboard(this.AdminContext(), &pb.ComposeAdminDashboardRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["dashboard"] = maps.Map{
		"countServers":      resp.CountServers,
		"countNodeClusters": resp.CountNodeClusters,
		"countNodes":        resp.CountNodes,
		"countUsers":        resp.CountUsers,
		"countAPINodes":     resp.CountAPINodes,
		"countDBNodes":      resp.CountDBNodes,
		"countUserNodes":    resp.CountUserNodes,

		"canGoServers":  configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeServer),
		"canGoNodes":    configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeNode),
		"canGoSettings": configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeSetting),
		"canGoUsers":    configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeUser),
	}

	// 今日流量
	todayTrafficBytes := int64(0)
	if len(resp.DailyTrafficStats) > 0 {
		todayTrafficBytes = resp.DailyTrafficStats[len(resp.DailyTrafficStats)-1].Bytes
	}
	todayTrafficString := numberutils.FormatBits(todayTrafficBytes * 8)
	result := regexp.MustCompile(`^(?U)(.+)([a-zA-Z]+)$`).FindStringSubmatch(todayTrafficString)
	if len(result) > 2 {
		this.Data["todayTraffic"] = result[1]
		this.Data["todayTrafficUnit"] = result[2]
	} else {
		this.Data["todayTraffic"] = todayTrafficString
		this.Data["todayTrafficUnit"] = ""
	}

	// 24小时流量趋势
	{
		statMaps := []maps.Map{}
		for _, stat := range resp.HourlyTrafficStats {
			statMaps = append(statMaps, maps.Map{
				"count": math.Ceil((float64(stat.Bytes)*8/1000/1000/1000)*1000) / 1000,
				"hour":  stat.Hour[8:],
			})
		}
		this.Data["hourlyTrafficStats"] = statMaps
	}

	// 15天流量趋势
	{
		statMaps := []maps.Map{}
		for _, stat := range resp.DailyTrafficStats {
			statMaps = append(statMaps, maps.Map{
				"count": math.Ceil((float64(stat.Bytes)*8/1000/1000/1000)*1000) / 1000,
				"day":   stat.Day[4:6] + "月" + stat.Day[6:] + "日",
			})
		}
		this.Data["dailyTrafficStats"] = statMaps
	}

	this.Show()
}
