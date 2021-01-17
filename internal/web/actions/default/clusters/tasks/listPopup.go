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
	resp, err := this.RPC().NodeTaskRPC().FindNodeClusterTasks(this.AdminContext(), &pb.FindNodeClusterTasksRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	countTasks := 0
	clusterMaps := []maps.Map{}
	for _, cluster := range resp.ClusterTasks {
		taskMaps := []maps.Map{}
		for _, task := range cluster.NodeTasks {
			countTasks++
			taskMaps = append(taskMaps, maps.Map{
				"id":   task.Id,
				"type": task.Type,
				"node": maps.Map{
					"id":   task.Node.Id,
					"name": task.Node.Name,
				},
				"isOk":        task.IsOk,
				"error":       task.Error,
				"isDone":      task.IsDone,
				"updatedTime": timeutil.FormatTime("Y-m-d H:i:s", task.UpdatedAt),
			})
		}

		clusterMaps = append(clusterMaps, maps.Map{
			"id":    cluster.ClusterId,
			"name":  cluster.ClusterName,
			"tasks": taskMaps,
		})
	}
	this.Data["clusters"] = clusterMaps
	this.Data["countTasks"] = countTasks
}
