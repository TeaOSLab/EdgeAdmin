package ipadminutils

import (
	"context"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/messageconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
)

// 通知使用此WAF策略的集群更新
func NotifyUpdateToClustersWithFirewallPolicyId(ctx context.Context, firewallPolicyId int64) error {
	client, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	resp, err := client.ServerRPC().FindAllEnabledServersWithHTTPFirewallPolicyId(ctx, &pb.FindAllEnabledServersWithHTTPFirewallPolicyIdRequest{FirewallPolicyId: firewallPolicyId})
	if err != nil {
		return err
	}
	clusterIds := []int64{}
	for _, server := range resp.Servers {
		if !lists.ContainsInt64(clusterIds, server.Cluster.Id) {
			clusterIds = append(clusterIds, server.Cluster.Id)
		}
	}
	for _, clusterId := range clusterIds {
		_, err = nodeutils.SendMessageToCluster(ctx, clusterId, messageconfigs.MessageCodeIPListChanged, &messageconfigs.IPListChangedMessage{}, 3)
		if err != nil {
			return err
		}
	}
	return nil
}
