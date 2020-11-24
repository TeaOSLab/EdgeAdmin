package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/messageconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/logs"
	"time"
)

func init() {
	events.On(events.EventStart, func() {
		task := NewSyncClusterTask()
		go task.Start()
	})
}

// 自动同步集群任务
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
	if !setup.IsConfigured() {
		return nil
	}

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	ctx := rpcClient.Context(0)
	resp, err := rpcClient.NodeClusterRPC().FindAllChangedNodeClusters(ctx, &pb.FindAllChangedNodeClustersRequest{})
	if err != nil {
		return err
	}

	for _, cluster := range resp.Clusters {
		_, err := rpcClient.NodeRPC().SyncNodesVersionWithCluster(ctx, &pb.SyncNodesVersionWithClusterRequest{
			ClusterId: cluster.Id,
		})
		if err != nil {
			return err
		}

		// 发送通知
		_, err = nodeutils.SendMessageToCluster(ctx, cluster.Id, messageconfigs.MessageCodeConfigChanged, &messageconfigs.ConfigChangedMessage{}, 10)
		if err != nil {
			return err
		}
	}
	return nil
}
