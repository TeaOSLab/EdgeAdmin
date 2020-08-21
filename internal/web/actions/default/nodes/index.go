package nodes

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs/nodes"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().NodeRPC().CountAllEnabledNodes(this.AdminContext(), &pb.CountAllEnabledNodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	page := this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	nodesResp, err := this.RPC().NodeRPC().ListEnabledNodes(this.AdminContext(), &pb.ListEnabledNodesRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	nodeMaps := []maps.Map{}
	for _, node := range nodesResp.Nodes {
		// 状态
		status := &nodes.NodeStatus{}
		if len(node.Status) > 0 && node.Status != "null" {
			err = json.Unmarshal([]byte(node.Status), &status)
			if err != nil {
				logs.Error(err)
				continue
			}
			status.IsActive = time.Now().Unix()-status.UpdatedAt < 120 // 2分钟之内认为活跃
		}

		nodeMaps = append(nodeMaps, maps.Map{
			"id":   node.Id,
			"name": node.Name,
			"status": maps.Map{
				"isActive":     status.IsActive,
				"updatedAt":    status.UpdatedAt,
				"hostname":     status.Hostname,
				"cpuUsage":     status.CPUUsage,
				"cpuUsageText": fmt.Sprintf("%.2f%%", status.CPUUsage*100),
				"memUsage":     status.MemoryUsage,
				"memUsageText": fmt.Sprintf("%.2f%%", status.MemoryUsage*100),
			},
			"cluster": maps.Map{
				"id":   node.Cluster.Id,
				"name": node.Cluster.Name,
			},
		})
	}
	this.Data["nodes"] = nodeMaps

	this.Show()
}
