package logs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	if this.ParamString("type") == "unread" {
		this.FirstMenu("unread")
	} else {
		this.FirstMenu("index")
	}
}

func (this *IndexAction) RunGet(params struct {
	DayFrom string
	DayTo   string
	Keyword string
	Level   string
	Type    string
}) {
	this.Data["dayFrom"] = params.DayFrom
	this.Data["dayTo"] = params.DayTo
	this.Data["keyword"] = params.Keyword
	this.Data["level"] = params.Level
	this.Data["type"] = params.Type

	// 未读数量
	countUnreadResp, err := this.RPC().NodeLogRPC().CountNodeLogs(this.AdminContext(), &pb.CountNodeLogsRequest{
		Role:     nodeconfigs.NodeRoleNode,
		IsUnread: true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countUnreadLogs"] = countUnreadResp.Count

	// 日志数量
	countResp, err := this.RPC().NodeLogRPC().CountNodeLogs(this.AdminContext(), &pb.CountNodeLogsRequest{
		NodeId:   0,
		Role:     nodeconfigs.NodeRoleNode,
		DayFrom:  params.DayFrom,
		DayTo:    params.DayTo,
		Keyword:  params.Keyword,
		Level:    params.Level,
		IsUnread: params.Type == "unread",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	logsResp, err := this.RPC().NodeLogRPC().ListNodeLogs(this.AdminContext(), &pb.ListNodeLogsRequest{
		NodeId:   0,
		Role:     nodeconfigs.NodeRoleNode,
		DayFrom:  params.DayFrom,
		DayTo:    params.DayTo,
		Keyword:  params.Keyword,
		Level:    params.Level,
		IsUnread: params.Type == "unread",
		Offset:   page.Offset,
		Size:     page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	logs := []maps.Map{}
	for _, log := range logsResp.NodeLogs {
		// 节点信息
		nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: log.NodeId})
		if err != nil {
			continue
		}
		node := nodeResp.Node
		if node == nil || node.NodeCluster == nil {
			continue
		}

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
			"id":          log.Id,
			"tag":         log.Tag,
			"description": log.Description,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt),
			"level":       log.Level,
			"isToday":     timeutil.FormatTime("Y-m-d", log.CreatedAt) == timeutil.Format("Y-m-d"),
			"count":       log.Count,
			"isRead":      log.IsRead,
			"node": maps.Map{
				"id": node.Id,
				"cluster": maps.Map{
					"id":   node.NodeCluster.Id,
					"name": node.NodeCluster.Name,
				},
				"name": node.Name,
			},
			"server": serverMap,
		})
	}
	this.Data["logs"] = logs

	this.Show()
}
