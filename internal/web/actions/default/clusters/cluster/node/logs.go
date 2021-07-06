package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
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
	NodeId  int64
	DayFrom string
	DayTo   string
	Keyword string
	Level   string
}) {
	// 初始化节点信息（用于菜单）
	err := nodeutils.InitNodeInfo(this, params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["nodeId"] = params.NodeId
	this.Data["dayFrom"] = params.DayFrom
	this.Data["dayTo"] = params.DayTo
	this.Data["keyword"] = params.Keyword
	this.Data["level"] = params.Level

	countResp, err := this.RPC().NodeLogRPC().CountNodeLogs(this.AdminContext(), &pb.CountNodeLogsRequest{
		Role:    "node",
		NodeId:  params.NodeId,
		DayFrom: params.DayFrom,
		DayTo:   params.DayTo,
		Keyword: params.Keyword,
		Level:   params.Level,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count, 20)

	logsResp, err := this.RPC().NodeLogRPC().ListNodeLogs(this.AdminContext(), &pb.ListNodeLogsRequest{
		NodeId:  params.NodeId,
		Role:    "node",
		DayFrom: params.DayFrom,
		DayTo:   params.DayTo,
		Keyword: params.Keyword,
		Level:   params.Level,
		Offset:  page.Offset,
		Size:    page.Size,
	})

	logs := []maps.Map{}
	for _, log := range logsResp.NodeLogs {
		logs = append(logs, maps.Map{
			"tag":         log.Tag,
			"description": log.Description,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt),
			"level":       log.Level,
			"isToday":     timeutil.FormatTime("Y-m-d", log.CreatedAt) == timeutil.Format("Y-m-d"),
			"count":       log.Count,
		})
	}
	this.Data["logs"] = logs

	this.Data["page"] = page.AsHTML()

	this.Show()
}
