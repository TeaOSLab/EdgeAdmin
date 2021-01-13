package log

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"strings"
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

	countResp, err := this.RPC().LogRPC().CountLogs(this.AdminContext(), &pb.CountLogRequest{
		DayFrom:  params.DayFrom,
		DayTo:    params.DayTo,
		Keyword:  params.Keyword,
		UserType: params.UserType,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	logsResp, err := this.RPC().LogRPC().ListLogs(this.AdminContext(), &pb.ListLogsRequest{
		Offset:   page.Offset,
		Size:     page.Size,
		DayFrom:  params.DayFrom,
		DayTo:    params.DayTo,
		Keyword:  params.Keyword,
		UserType: params.UserType,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	logMaps := []maps.Map{}
	for _, log := range logsResp.Logs {
		regionName := ""
		regionResp, err := this.RPC().IPLibraryRPC().LookupIPRegion(this.AdminContext(), &pb.LookupIPRegionRequest{Ip: log.Ip})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if regionResp.IpRegion != nil {
			pieces := []string{regionResp.IpRegion.Summary}
			if len(regionResp.IpRegion.Isp) > 0 && !lists.ContainsString(pieces, regionResp.IpRegion.Isp) {
				pieces = append(pieces, "| "+regionResp.IpRegion.Isp)
			}
			regionName = strings.Join(pieces, " ")
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
