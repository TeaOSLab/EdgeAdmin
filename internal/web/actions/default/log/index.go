package log

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("log", "log", "list")
}

func (this *IndexAction) RunGet(params struct {
	DayFrom  string
	DayTo    string
	Keyword  string
	UserType string
	Level    string
}) {
	// 读取配置
	config, err := configloaders.LoadLogConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["logConfig"] = config

	this.Data["dayFrom"] = params.DayFrom
	this.Data["dayTo"] = params.DayTo
	this.Data["keyword"] = params.Keyword
	this.Data["userType"] = params.UserType

	// 级别
	this.Data["level"] = params.Level
	this.Data["levelOptions"] = []maps.Map{
		{
			"code": "info",
			"name": this.Lang(codes.Level_Info),
		},
		{
			"code": "warn",
			"name": this.Lang(codes.Level_Warn),
		},
		{
			"code": "error",
			"name": this.Lang(codes.Level_Error),
		},
	}

	countResp, err := this.RPC().LogRPC().CountLogs(this.AdminContext(), &pb.CountLogRequest{
		DayFrom:  params.DayFrom,
		DayTo:    params.DayTo,
		Keyword:  params.Keyword,
		UserType: params.UserType,
		Level:    params.Level,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var page = this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	logsResp, err := this.RPC().LogRPC().ListLogs(this.AdminContext(), &pb.ListLogsRequest{
		Offset:   page.Offset,
		Size:     page.Size,
		DayFrom:  params.DayFrom,
		DayTo:    params.DayTo,
		Keyword:  params.Keyword,
		UserType: params.UserType,
		Level:    params.Level,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var logMaps = []maps.Map{}
	for _, log := range logsResp.Logs {
		var regionName = ""
		var ipRegion = iplibrary.LookupIP(log.Ip)
		if ipRegion != nil && ipRegion.IsOk() {
			regionName = ipRegion.Summary()
		}

		logMaps = append(logMaps, maps.Map{
			"id":          log.Id,
			"adminId":     log.AdminId,
			"userId":      log.UserId,
			"description": log.Description,
			"userName":    log.UserName,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt),
			"level":       log.Level,
			"type":        log.Type,
			"ip":          log.Ip,
			"region":      regionName,
			"action":      log.Action,
		})
	}
	this.Data["logs"] = logMaps

	this.Show()
}
