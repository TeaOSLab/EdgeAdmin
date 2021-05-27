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
}

func (this *IndexAction) RunGet(params struct {
	DayFrom string
	DayTo   string
	Keyword string
	Level   string
}) {
	this.Data["dayFrom"] = params.DayFrom
	this.Data["dayTo"] = params.DayTo
	this.Data["keyword"] = params.Keyword
	this.Data["level"] = params.Level

	countResp, err := this.RPC().NodeLogRPC().CountNodeLogs(this.AdminContext(), &pb.CountNodeLogsRequest{
		NodeId:  0,
		Role:    nodeconfigs.NodeRoleDNS,
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
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	logsResp, err := this.RPC().NodeLogRPC().ListNodeLogs(this.AdminContext(), &pb.ListNodeLogsRequest{
		NodeId:  0,
		Role:    nodeconfigs.NodeRoleDNS,
		DayFrom: params.DayFrom,
		DayTo:   params.DayTo,
		Keyword: params.Keyword,
		Level:   params.Level,
		Offset:  page.Offset,
		Size:    page.Size,
	})

	logs := []maps.Map{}
	for _, log := range logsResp.NodeLogs {
		// 节点信息
		nodeResp, err := this.RPC().NSNodeRPC().FindEnabledNSNode(this.AdminContext(), &pb.FindEnabledNSNodeRequest{NsNodeId: log.NodeId})
		if err != nil {
			continue
		}
		node := nodeResp.NsNode
		if node == nil || node.NsCluster == nil {
			continue
		}

		logs = append(logs, maps.Map{
			"tag":         log.Tag,
			"description": log.Description,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt),
			"level":       log.Level,
			"isToday":     timeutil.FormatTime("Y-m-d", log.CreatedAt) == timeutil.Format("Y-m-d"),
			"count":       log.Count,
			"node": maps.Map{
				"id": node.Id,
				"cluster": maps.Map{
					"id":   node.NsCluster.Id,
					"name": node.NsCluster.Name,
				},
				"name": node.Name,
			},
		})
	}
	this.Data["logs"] = logs

	this.Show()
}
