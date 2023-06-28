package logs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/nodelogutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	var paramType = this.ParamString("type")
	if paramType == "unread" {
		this.FirstMenu("unread")
	} else if paramType == "needFix" {
		this.FirstMenu("needFix")
	} else {
		this.FirstMenu("index")
	}
}

func (this *IndexAction) RunGet(params struct {
	DayFrom   string
	DayTo     string
	Keyword   string
	Level     string
	Type      string // unread, needFix
	Tag       string
	ClusterId int64
	NodeId    int64
}) {
	this.Data["dayFrom"] = params.DayFrom
	this.Data["dayTo"] = params.DayTo
	this.Data["keyword"] = params.Keyword
	this.Data["searchedKeyword"] = params.Keyword
	this.Data["level"] = params.Level
	this.Data["type"] = params.Type
	this.Data["tag"] = params.Tag
	this.Data["clusterId"] = params.ClusterId
	this.Data["nodeId"] = params.NodeId

	var fixedState configutils.BoolState = 0
	var allServers = false
	if params.Type == "needFix" {
		fixedState = configutils.BoolStateNo
		allServers = true
	}

	// 常见标签
	this.Data["tags"] = nodelogutils.FindNodeCommonTags(this.LangCode())

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

	// 需要修复数量
	countNeedFixResp, err := this.RPC().NodeLogRPC().CountNodeLogs(this.AdminContext(), &pb.CountNodeLogsRequest{
		Role:       nodeconfigs.NodeRoleNode,
		AllServers: true,
		FixedState: int32(configutils.BoolStateNo),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countNeedFixLogs"] = countNeedFixResp.Count

	// 日志数量
	countResp, err := this.RPC().NodeLogRPC().CountNodeLogs(this.AdminContext(), &pb.CountNodeLogsRequest{
		NodeClusterId: params.ClusterId,
		NodeId:        params.NodeId,
		Role:          nodeconfigs.NodeRoleNode,
		DayFrom:       params.DayFrom,
		DayTo:         params.DayTo,
		Keyword:       params.Keyword,
		Level:         params.Level,
		IsUnread:      params.Type == "unread",
		Tag:           params.Tag,
		FixedState:    int32(fixedState),
		AllServers:    allServers,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	this.Data["countLogs"] = count

	var page = this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	logsResp, err := this.RPC().NodeLogRPC().ListNodeLogs(this.AdminContext(), &pb.ListNodeLogsRequest{
		NodeClusterId: params.ClusterId,
		NodeId:        params.NodeId,
		Role:          nodeconfigs.NodeRoleNode,
		DayFrom:       params.DayFrom,
		DayTo:         params.DayTo,
		Keyword:       params.Keyword,
		Level:         params.Level,
		IsUnread:      params.Type == "unread",
		Tag:           params.Tag,
		FixedState:    int32(fixedState),
		AllServers:    allServers,
		Offset:        page.Offset,
		Size:          page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var firstUnreadNodeMap maps.Map = nil

	var logs = []maps.Map{}
	for _, log := range logsResp.NodeLogs {
		// 节点信息
		nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: log.NodeId})
		if err != nil {
			continue
		}
		var node = nodeResp.Node
		if node == nil || node.NodeCluster == nil {
			continue
		}

		if params.Type == "unread" && firstUnreadNodeMap == nil {
			firstUnreadNodeMap = maps.Map{
				"id":   node.Id,
				"name": node.Name,
			}
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

		var isFixed = true
		if !log.IsFixed && log.ServerId > 0 && lists.ContainsString([]string{"success", "warning", "error"}, log.Level) {
			isFixed = false
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
			"isFixed":     isFixed,
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

	this.Data["firstUnreadNode"] = firstUnreadNodeMap

	this.Show()
}
