package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type ListPopupAction struct {
	actionutils.ParentAction
}

func (this *ListPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *ListPopupAction) RunGet(params struct{}) {
	this.retrieveTasks()

	this.Show()
}

func (this *ListPopupAction) RunPost(params struct {
	Must *actions.Must
}) {
	this.retrieveTasks()
	this.Success()
}

func (this *ListPopupAction) retrieveTasks() {
	resp, err := this.RPC().DNSTaskRPC().FindAllDoingDNSTasks(this.AdminContext(), &pb.FindAllDoingDNSTasksRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	taskMaps := []maps.Map{}
	for _, task := range resp.DnsTasks {
		var clusterMap maps.Map = nil
		var nodeMap maps.Map = nil
		var serverMap maps.Map = nil
		var domainMap maps.Map = nil

		if task.NodeCluster != nil {
			clusterMap = maps.Map{
				"id":   task.NodeCluster.Id,
				"name": task.NodeCluster.Name,
			}
		}
		if task.Node != nil {
			nodeMap = maps.Map{
				"id":   task.Node.Id,
				"name": task.Node.Name,
			}
		}
		if task.Server != nil {
			serverMap = maps.Map{
				"id":   task.Server.Id,
				"name": task.Server.Name,
			}
		}
		if task.DnsDomain != nil {
			domainMap = maps.Map{
				"id":   task.DnsDomain.Id,
				"name": task.DnsDomain.Name,
			}
		}

		taskMaps = append(taskMaps, maps.Map{
			"id":          task.Id,
			"type":        task.Type,
			"isDone":      task.IsDone,
			"isOk":        task.IsOk,
			"error":       task.Error,
			"updatedTime": timeutil.FormatTime("Y-m-d H:i:s", task.UpdatedAt),
			"cluster":     clusterMap,
			"node":        nodeMap,
			"server":      serverMap,
			"domain":      domainMap,
		})
	}
	this.Data["tasks"] = taskMaps
}
