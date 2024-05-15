package nodeutils

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"sort"
	"strconv"
	"sync"
)

// MessageResult 和节点消息通讯结果定义
type MessageResult struct {
	NodeId   int64  `json:"nodeId"`
	NodeName string `json:"nodeName"`
	IsOK     bool   `json:"isOk"`
	Message  string `json:"message"`
}

// SendMessageToCluster 向集群发送命令消息
func SendMessageToCluster(ctx context.Context, clusterId int64, code string, msg any, timeoutSeconds int32, availableNodesOnly bool) (results []*MessageResult, err error) {
	results = []*MessageResult{}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return results, err
	}

	defaultRPCClient, err := rpc.SharedRPC()
	if err != nil {
		return results, err
	}

	// 获取所有节点
	nodesResp, err := defaultRPCClient.NodeRPC().FindAllEnabledNodesWithNodeClusterId(ctx, &pb.FindAllEnabledNodesWithNodeClusterIdRequest{
		NodeClusterId:    clusterId,
		IncludeSecondary: true,
	})
	if err != nil {
		return results, err
	}
	var nodes = nodesResp.Nodes
	if len(nodes) == 0 {
		return results, nil
	}

	if availableNodesOnly {
		var newNodes []*pb.Node
		for _, node := range nodes {
			if !node.IsOn {
				continue
			}
			newNodes = append(newNodes, node)
		}
		nodes = newNodes
	}

	var rpcMap = map[int64]*rpc.RPCClient{} // apiNodeId => RPCClient
	var locker = &sync.Mutex{}

	var wg = &sync.WaitGroup{}
	wg.Add(len(nodes))

	for _, node := range nodes {
		// TODO 检查是否在线

		if !node.IsOn {
			locker.Lock()
			results = append(results, &MessageResult{
				NodeId:   node.Id,
				NodeName: node.Name,
				IsOK:     false,
				Message:  "节点尚未启用",
			})
			locker.Unlock()
			wg.Done()
			continue
		}

		if len(node.ConnectedAPINodeIds) == 0 {
			locker.Lock()
			results = append(results, &MessageResult{
				NodeId:   node.Id,
				NodeName: node.Name,
				IsOK:     false,
				Message:  "节点尚未连接到API",
			})
			locker.Unlock()
			wg.Done()
			continue
		}

		// 获取API节点信息
		apiNodeId := node.ConnectedAPINodeIds[0]
		rpcClient, ok := rpcMap[apiNodeId]
		if !ok {
			apiNodeResp, err := defaultRPCClient.APINodeRPC().FindEnabledAPINode(ctx, &pb.FindEnabledAPINodeRequest{ApiNodeId: apiNodeId})
			if err != nil {
				locker.Lock()
				results = append(results, &MessageResult{
					NodeId:   node.Id,
					NodeName: node.Name,
					IsOK:     false,
					Message:  "无法读取对应的API节点信息：" + err.Error(),
				})
				locker.Unlock()
				wg.Done()
				continue
			}

			if apiNodeResp.ApiNode == nil {
				locker.Lock()
				results = append(results, &MessageResult{
					NodeId:   node.Id,
					NodeName: node.Name,
					IsOK:     false,
					Message:  "无法读取对应的API节点信息：API节点ID：" + strconv.FormatInt(apiNodeId, 10),
				})
				locker.Unlock()
				wg.Done()
				continue
			}
			apiNode := apiNodeResp.ApiNode

			apiRPCClient, err := rpc.NewRPCClient(&configs.APIConfig{
				RPCEndpoints: apiNode.AccessAddrs,
				NodeId:       apiNode.UniqueId,
				Secret:       apiNode.Secret,
			}, false)
			if err != nil {
				locker.Lock()
				results = append(results, &MessageResult{
					NodeId:   node.Id,
					NodeName: node.Name,
					IsOK:     false,
					Message:  "初始化API节点错误：API节点ID：" + strconv.FormatInt(apiNodeId, 10) + "：" + err.Error(),
				})
				locker.Unlock()
				wg.Done()
				continue
			}
			rpcMap[apiNodeId] = apiRPCClient
			rpcClient = apiRPCClient
		}

		// 发送消息
		go func(node *pb.Node) {
			defer wg.Done()

			result, err := rpcClient.NodeRPC().SendCommandToNode(ctx, &pb.NodeStreamMessage{
				NodeId:         node.Id,
				TimeoutSeconds: timeoutSeconds,
				Code:           code,
				DataJSON:       msgJSON,
			})
			if err != nil {
				locker.Lock()
				results = append(results, &MessageResult{
					NodeId:   node.Id,
					NodeName: node.Name,
					IsOK:     false,
					Message:  "API返回错误：" + err.Error(),
				})
				locker.Unlock()
				return
			}

			locker.Lock()
			results = append(results, &MessageResult{
				NodeId:   node.Id,
				NodeName: node.Name,
				IsOK:     result.IsOk,
				Message:  result.Message,
			})
			locker.Unlock()
		}(node)
	}
	wg.Wait()

	// 对结果进行排序
	if len(results) > 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i].NodeId < results[j].NodeId
		})
	}

	// 关闭RPC
	for _, rpcClient := range rpcMap {
		_ = rpcClient.Close()
	}

	return
}

// SendMessageToNodeIds 向一组节点发送命令消息
func SendMessageToNodeIds(ctx context.Context, nodeIds []int64, code string, msg interface{}, timeoutSeconds int32) (results []*MessageResult, err error) {
	results = []*MessageResult{}
	if len(nodeIds) == 0 {
		return
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return results, err
	}

	defaultRPCClient, err := rpc.SharedRPC()
	if err != nil {
		return results, err
	}

	// 获取所有节点
	nodesResp, err := defaultRPCClient.NodeRPC().FindEnabledNodesWithIds(ctx, &pb.FindEnabledNodesWithIdsRequest{NodeIds: nodeIds})
	if err != nil {
		return nil, err
	}
	nodes := nodesResp.Nodes
	if len(nodes) == 0 {
		return results, nil
	}

	rpcMap := map[int64]*rpc.RPCClient{} // apiNodeId => RPCClient
	locker := &sync.Mutex{}

	wg := &sync.WaitGroup{}
	wg.Add(len(nodes))
	for _, node := range nodes {
		if !node.IsActive {
			locker.Lock()
			results = append(results, &MessageResult{
				NodeId:   node.Id,
				NodeName: node.Name,
				IsOK:     false,
				Message:  "节点不在线",
			})
			locker.Unlock()
			wg.Done()
			continue
		}

		if !node.IsOn {
			if !node.IsActive {
				locker.Lock()
				results = append(results, &MessageResult{
					NodeId:   node.Id,
					NodeName: node.Name,
					IsOK:     false,
					Message:  "节点未启用",
				})
				locker.Unlock()
				wg.Done()
				continue
			}
		}

		if len(node.ConnectedAPINodeIds) == 0 {
			locker.Lock()
			results = append(results, &MessageResult{
				NodeId:   node.Id,
				NodeName: node.Name,
				IsOK:     false,
				Message:  "节点尚未连接到API",
			})
			locker.Unlock()
			wg.Done()
			continue
		}

		// 获取API节点信息
		apiNodeId := node.ConnectedAPINodeIds[0]
		rpcClient, ok := rpcMap[apiNodeId]
		if !ok {
			apiNodeResp, err := defaultRPCClient.APINodeRPC().FindEnabledAPINode(ctx, &pb.FindEnabledAPINodeRequest{ApiNodeId: apiNodeId})
			if err != nil {
				locker.Lock()
				results = append(results, &MessageResult{
					NodeId:   node.Id,
					NodeName: node.Name,
					IsOK:     false,
					Message:  "无法读取对应的API节点信息：" + err.Error(),
				})
				locker.Unlock()
				wg.Done()
				continue
			}

			if apiNodeResp.ApiNode == nil {
				locker.Lock()
				results = append(results, &MessageResult{
					NodeId:   node.Id,
					NodeName: node.Name,
					IsOK:     false,
					Message:  "无法读取对应的API节点信息：API节点ID：" + strconv.FormatInt(apiNodeId, 10),
				})
				locker.Unlock()
				wg.Done()
				continue
			}
			apiNode := apiNodeResp.ApiNode

			apiRPCClient, err := rpc.NewRPCClient(&configs.APIConfig{
				RPCEndpoints: apiNode.AccessAddrs,
				NodeId:       apiNode.UniqueId,
				Secret:       apiNode.Secret,
			}, false)
			if err != nil {
				locker.Lock()
				results = append(results, &MessageResult{
					NodeId:   node.Id,
					NodeName: node.Name,
					IsOK:     false,
					Message:  "初始化API节点错误：API节点ID：" + strconv.FormatInt(apiNodeId, 10) + "：" + err.Error(),
				})
				locker.Unlock()
				wg.Done()
				continue
			}
			rpcMap[apiNodeId] = apiRPCClient
			rpcClient = apiRPCClient
		}

		// 发送消息
		go func(node *pb.Node) {
			defer wg.Done()

			result, err := rpcClient.NodeRPC().SendCommandToNode(ctx, &pb.NodeStreamMessage{
				NodeId:         node.Id,
				TimeoutSeconds: timeoutSeconds,
				Code:           code,
				DataJSON:       msgJSON,
			})
			if err != nil {
				locker.Lock()
				results = append(results, &MessageResult{
					NodeId:   node.Id,
					NodeName: node.Name,
					IsOK:     false,
					Message:  "API返回错误：" + err.Error(),
				})
				locker.Unlock()
				return
			}

			locker.Lock()
			results = append(results, &MessageResult{
				NodeId:   node.Id,
				NodeName: node.Name,
				IsOK:     result.IsOk,
				Message:  result.Message,
			})
			locker.Unlock()
		}(node)
	}
	wg.Wait()

	// 对结果进行排序
	if len(results) > 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i].NodeId < results[j].NodeId
		})
	}

	// 关闭RPC
	for _, rpcClient := range rpcMap {
		_ = rpcClient.Close()
	}

	return
}
