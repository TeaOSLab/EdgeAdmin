package tasks

import (
	"context"
	"crypto/tls"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/goman"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

func init() {
	events.On(events.EventStart, func() {
		task := NewSyncAPINodesTask()
		goman.New(func() {
			task.Start()
		})
	})
}

// SyncAPINodesTask API节点同步任务
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

	config, err := configs.LoadAPIConfig()
	if err != nil {
		return err
	}

	// 是否禁止自动升级
	if config.RPCDisableUpdate {
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

	var newEndpoints = []string{}
	for _, node := range resp.ApiNodes {
		if !node.IsOn {
			continue
		}
		newEndpoints = append(newEndpoints, node.AccessAddrs...)
	}

	// 和现有的对比
	if this.isSame(newEndpoints, config.RPCEndpoints) {
		return nil
	}

	// 测试是否有API节点可用
	hasOk := this.testEndpoints(newEndpoints)
	if !hasOk {
		return nil
	}

	// 修改RPC对象配置
	config.RPCEndpoints = newEndpoints
	err = rpcClient.UpdateConfig(config)
	if err != nil {
		return err
	}

	// 保存到文件
	err = config.WriteFile(Tea.ConfigFile(configs.ConfigFileName))
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

func (this *SyncAPINodesTask) testEndpoints(endpoints []string) bool {
	if len(endpoints) == 0 {
		return false
	}

	var wg = sync.WaitGroup{}
	wg.Add(len(endpoints))

	var ok = false

	for _, endpoint := range endpoints {
		go func(endpoint string) {
			defer wg.Done()

			u, err := url.Parse(endpoint)
			if err != nil {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer func() {
				cancel()
			}()
			var conn *grpc.ClientConn

			if u.Scheme == "http" {
				conn, err = grpc.DialContext(ctx, u.Host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
			} else if u.Scheme == "https" {
				conn, err = grpc.DialContext(ctx, u.Host, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: true,
				})), grpc.WithBlock())
			}
			if err != nil {
				return
			}
			_ = conn.Close()

			ok = true
		}(endpoint)
	}
	wg.Wait()

	return ok
}
