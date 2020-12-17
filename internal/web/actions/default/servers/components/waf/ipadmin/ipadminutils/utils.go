package ipadminutils

import (
	"context"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/messageconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

// 通知使用此WAF策略的集群更新
func NotifyUpdateToClustersWithFirewallPolicyId(ctx context.Context, firewallPolicyId int64) error {
	client, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	resp, err := client.NodeClusterRPC().FindAllEnabledNodeClustersWithHTTPFirewallPolicyId(ctx, &pb.FindAllEnabledNodeClustersWithHTTPFirewallPolicyIdRequest{HttpFirewallPolicyId: firewallPolicyId})
	if err != nil {
		return err
	}
	for _, cluster := range resp.NodeClusters {
		_, err = nodeutils.SendMessageToCluster(ctx, cluster.Id, messageconfigs.MessageCodeIPListChanged, &messageconfigs.IPListChangedMessage{}, 3)
		if err != nil {
			return err
		}
	}
	return nil
}
