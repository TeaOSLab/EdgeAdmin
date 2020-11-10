package log

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("log", "log", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().LogRPC().CountLogs(this.AdminContext(), &pb.CountLogRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	logsResp, err := this.RPC().LogRPC().ListLogs(this.AdminContext(), &pb.ListLogsRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	logMaps := []maps.Map{}
	for _, log := range logsResp.Logs {
		logMaps = append(logMaps, maps.Map{
			"description": log.Description,
			"userName":    log.Description,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt),
			"level":       log.Level,
			"type":        log.Type,
			"ip":          log.Ip,
		})
	}
	this.Data["logs"] = logMaps

	this.Show()
}
