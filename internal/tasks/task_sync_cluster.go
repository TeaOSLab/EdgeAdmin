package tasks

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/goman"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/messageconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"time"
)

func init() {
	events.On(events.EventStart, func() {
		task := NewSyncClusterTask()
		goman.New(func() {
			task.Start()
		})
	})
}

// SyncClusterTask 自动同步集群任务
type SyncClusterTask struct {
}

func NewSyncClusterTask() *SyncClusterTask {
	return &SyncClusterTask{}
}

func (this *SyncClusterTask) Start() {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		err := this.loop()
		if err != nil {
			logs.Println("[TASK][SYNC_CLUSTER]" + err.Error())
		}
	}
}

func (this *SyncClusterTask) loop() error {
	// 如果还没有安装直接返回
	if !setup.IsConfigured() || teaconst.IsRecoverMode {
		return nil
	}

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	ctx := rpcClient.Context(0)

	tasksResp, err := rpcClient.NodeTaskRPC().FindNotifyingNodeTasks(ctx, &pb.FindNotifyingNodeTasksRequest{Size: 500})
	if err != nil {
		return err
	}
	nodeIds := []int64{}
	taskIds := []int64{}
	for _, task := range tasksResp.NodeTasks {
		if !lists.ContainsInt64(nodeIds, task.Node.Id) {
			nodeIds = append(nodeIds, task.Node.Id)
		}
		taskIds = append(taskIds, task.Id)
	}
	if len(nodeIds) == 0 {
		return nil
	}
	_, err = nodeutils.SendMessageToNodeIds(ctx, nodeIds, messageconfigs.MessageCodeNewNodeTask, &messageconfigs.NewNodeTaskMessage{}, 3)
	if err != nil {
		return err
	}

	// 设置已通知
	_, err = rpcClient.NodeTaskRPC().UpdateNodeTasksNotified(ctx, &pb.UpdateNodeTasksNotifiedRequest{NodeTaskIds: taskIds})
	if err != nil {
		return err
	}

	return nil
}
