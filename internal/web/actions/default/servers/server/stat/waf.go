package stat

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type WafAction struct {
	actionutils.ParentAction
}

func (this *WafAction) Init() {
	this.Nav("", "stat", "")
	this.SecondMenu("waf")
}

func (this *WafAction) RunGet(params struct {
	ServerId int64
}) {
	// 统计数据
	resp, err := this.RPC().ServerHTTPFirewallDailyStatRPC().ComposeServerHTTPFirewallDashboard(this.AdminContext(), &pb.ComposeServerHTTPFirewallDashboardRequest{
		Day:      timeutil.Format("Ymd"),
		ServerId: params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["countDailyLog"] = resp.CountDailyLog
	this.Data["countDailyBlock"] = resp.CountDailyBlock
	this.Data["countDailyCaptcha"] = resp.CountDailyCaptcha
	this.Data["countWeeklyBlock"] = resp.CountWeeklyBlock
	this.Data["countMonthlyBlock"] = resp.CountMonthlyBlock

	// 分组
	groupStatMaps := []maps.Map{}
	for _, group := range resp.HttpFirewallRuleGroups {
		groupStatMaps = append(groupStatMaps, maps.Map{
			"group": maps.Map{
				"id":   group.HttpFirewallRuleGroup.Id,
				"name": group.HttpFirewallRuleGroup.Name,
			},
			"count": group.Count,
		})
	}
	this.Data["groupStats"] = groupStatMaps

	// 每日趋势
	logStatMaps := []maps.Map{}
	blockStatMaps := []maps.Map{}
	captchaStatMaps := []maps.Map{}
	for _, stat := range resp.LogDailyStats {
		logStatMaps = append(logStatMaps, maps.Map{
			"day":   stat.Day,
			"count": stat.Count,
		})
	}
	for _, stat := range resp.BlockDailyStats {
		blockStatMaps = append(blockStatMaps, maps.Map{
			"day":   stat.Day,
			"count": stat.Count,
		})
	}
	for _, stat := range resp.CaptchaDailyStats {
		captchaStatMaps = append(captchaStatMaps, maps.Map{
			"day":   stat.Day,
			"count": stat.Count,
		})
	}
	this.Data["logDailyStats"] = logStatMaps
	this.Data["blockDailyStats"] = blockStatMaps
	this.Data["captchaDailyStats"] = captchaStatMaps

	this.Show()
}
