package logs

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
	this.Nav("", "", "log")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().MessageTaskLogRPC().CountMessageTaskLogs(this.AdminContext(), &pb.CountMessageTaskLogsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	page := this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	logsResp, err := this.RPC().MessageTaskLogRPC().ListMessageTaskLogs(this.AdminContext(), &pb.ListMessageTaskLogsRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	logMaps := []maps.Map{}
	for _, log := range logsResp.MessageTaskLogs {
		if log.MessageTask.MessageRecipient != nil {
			log.MessageTask.User = log.MessageTask.MessageRecipient.User
		}
		logMaps = append(logMaps, maps.Map{
			"task": maps.Map{
				"id":      log.MessageTask.Id,
				"user":    log.MessageTask.User,
				"subject": log.MessageTask.Subject,
				"body":    log.MessageTask.Body,
				"instance": maps.Map{
					"id":   log.MessageTask.MessageMediaInstance.Id,
					"name": log.MessageTask.MessageMediaInstance.Name,
				},
			},
			"isOk":        log.IsOk,
			"error":       log.Error,
			"response":    log.Response,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt),
		})
	}
	this.Data["logs"] = logMaps

	this.Show()
}
