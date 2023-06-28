package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/nodelogutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
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
	Tag     string
}) {
	// 初始化节点信息（用于菜单）
	_, err := nodeutils.InitNodeInfo(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["tags"] = nodelogutils.FindNodeCommonTags(this.LangCode())

	this.Data["nodeId"] = params.NodeId
	this.Data["dayFrom"] = params.DayFrom
	this.Data["dayTo"] = params.DayTo
	this.Data["keyword"] = params.Keyword
	this.Data["level"] = params.Level
	this.Data["tag"] = params.Tag

	countResp, err := this.RPC().NodeLogRPC().CountNodeLogs(this.AdminContext(), &pb.CountNodeLogsRequest{
		Role:    nodeconfigs.NodeRoleNode,
		NodeId:  params.NodeId,
		DayFrom: params.DayFrom,
		DayTo:   params.DayTo,
		Keyword: params.Keyword,
		Level:   params.Level,
		Tag:     params.Tag,
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
		Tag:     params.Tag,
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	logs := []maps.Map{}
	for _, log := range logsResp.NodeLogs {
		// 服务信息
		var serverMap = maps.Map{"id": 0}
		if log.ServerId > 0 {
			serverResp, err := this.RPC().ServerRPC().FindEnabledUserServerBasic(this.AdminContext(), &pb.FindEnabledUserServerBasicRequest{ServerId: log.ServerId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			var server = serverResp.Server
			if server != nil {
				serverMap = maps.Map{"id": server.Id, "name": server.Name}
			}
		}

		logs = append(logs, maps.Map{
			"tag":         log.Tag,
			"description": log.Description,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt),
			"level":       log.Level,
			"isToday":     timeutil.FormatTime("Y-m-d", log.CreatedAt) == timeutil.Format("Y-m-d"),
			"count":       log.Count,
			"server":      serverMap,
		})
	}
	this.Data["logs"] = logs

	this.Data["page"] = page.AsHTML()

	this.Show()
}
