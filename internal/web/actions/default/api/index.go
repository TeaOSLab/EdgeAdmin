package api

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
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
	countResp, err := this.RPC().APINodeRPC().CountAllEnabledAPINodes(this.AdminContext(), &pb.CountAllEnabledAPINodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	nodeMaps := []maps.Map{}
	if count > 0 {
		nodesResp, err := this.RPC().APINodeRPC().ListEnabledAPINodes(this.AdminContext(), &pb.ListEnabledAPINodesRequest{
			Offset: page.Offset,
			Size:   page.Size,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		for _, node := range nodesResp.Nodes {
			// 状态
			status := &nodeconfigs.NodeStatus{}
			if len(node.StatusJSON) > 0 {
				err = json.Unmarshal(node.StatusJSON, &status)
				if err != nil {
					logs.Error(err)
					continue
				}
				status.IsActive = status.IsActive && time.Now().Unix()-status.UpdatedAt <= 60 // N秒之内认为活跃
			}

			nodeMaps = append(nodeMaps, maps.Map{
				"id":          node.Id,
				"isOn":        node.IsOn,
				"name":        node.Name,
				"accessAddrs": node.AccessAddrs,
				"status": maps.Map{
					"isActive":     status.IsActive,
					"updatedAt":    status.UpdatedAt,
					"hostname":     status.Hostname,
					"cpuUsage":     status.CPUUsage,
					"cpuUsageText": fmt.Sprintf("%.2f%%", status.CPUUsage*100),
					"memUsage":     status.MemoryUsage,
					"memUsageText": fmt.Sprintf("%.2f%%", status.MemoryUsage*100),
					"buildVersion": status.BuildVersion,
				},
			})
		}
	}
	this.Data["nodes"] = nodeMaps

	this.Show()
}
