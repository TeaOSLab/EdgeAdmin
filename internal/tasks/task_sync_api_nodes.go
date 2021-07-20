package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"sort"
	"strings"
	"time"
)

func init() {
	events.On(events.EventStart, func() {
		task := NewSyncAPINodesTask()
		go task.Start()
	})
}

// API节点同步任务
type SyncAPINodesTask struct {
}

func NewSyncAPINodesTask() *SyncAPINodesTask {
	return &SyncAPINodesTask{}
}

func (this *SyncAPINodesTask) Start() {
	ticker := time.NewTicker(5 * time.Minute)
	if Tea.IsTesting() {
		// 快速测试
		ticker = time.NewTicker(1 * time.Minute)
	}
	for range ticker.C {
		err := this.Loop()
		if err != nil {
			logs.Println("[TASK][SYNC_API_NODES]" + err.Error())
		}
	}
}

func (this *SyncAPINodesTask) Loop() error {
	// 如果还没有安装直接返回
	if !setup.IsConfigured() || teaconst.IsRecoverMode {
		return nil
	}

	// 获取所有可用的节点
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	resp, err := rpcClient.APINodeRPC().FindAllEnabledAPINodes(rpcClient.Context(0), &pb.FindAllEnabledAPINodesRequest{})
	if err != nil {
		return err
	}

	newEndpoints := []string{}
	for _, node := range resp.Nodes {
		if !node.IsOn {
			continue
		}
		newEndpoints = append(newEndpoints, node.AccessAddrs...)
	}

	// 和现有的对比
	config, err := configs.LoadAPIConfig()
	if err != nil {
		return err
	}
	if this.isSame(newEndpoints, config.RPC.Endpoints) {
		return nil
	}

	// 修改RPC对象配置
	config.RPC.Endpoints = newEndpoints
	err = rpcClient.UpdateConfig(config)
	if err != nil {
		return err
	}

	// 保存到文件
	err = config.WriteFile(Tea.ConfigFile("api.yaml"))
	if err != nil {
		return err
	}

	return nil
}

func (this *SyncAPINodesTask) isSame(endpoints1 []string, endpoints2 []string) bool {
	sort.Strings(endpoints1)
	sort.Strings(endpoints2)
	return strings.Join(endpoints1, "&") == strings.Join(endpoints2, "&")
}
