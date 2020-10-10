package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type LogsAction struct {
	actionutils.ParentAction
}

func (this *LogsAction) Init() {
	this.Nav("", "node", "log")
	this.SecondMenu("nodes")
}

func (this *LogsAction) RunGet(params struct {
	NodeId int64
}) {
	this.Data["nodeId"] = params.NodeId

	countResp, err := this.RPC().NodeLogRPC().CountNodeLogs(this.AdminContext(), &pb.CountNodeLogsRequest{
		Role:   "node",
		NodeId: params.NodeId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count, 20)

	logsResp, err := this.RPC().NodeLogRPC().ListNodeLogs(this.AdminContext(), &pb.ListNodeLogsRequest{
		NodeId: params.NodeId,
		Role:   "node",
		Offset: page.Offset,
		Size:   page.Size,
	})

	logs := []maps.Map{}
	for _, log := range logsResp.NodeLogs {
		logs = append(logs, maps.Map{
			"tag":         log.Tag,
			"description": log.Description,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt),
			"level":       log.Level,
			"isToday":     timeutil.FormatTime("Y-m-d", log.CreatedAt) == timeutil.Format("Y-m-d"),
		})
	}
	this.Data["logs"] = logs

	this.Data["page"] = page.AsHTML()

	this.Show()
}
