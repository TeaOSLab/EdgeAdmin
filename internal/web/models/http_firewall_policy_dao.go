package models

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
)

var SharedHTTPFirewallPolicyDAO = new(HTTPFirewallPolicyDAO)

type HTTPFirewallPolicyDAO struct {
}

// 查找缓存策略基本信息
func (this *HTTPFirewallPolicyDAO) FindEnabledPolicy(ctx context.Context, policyId int64) (*pb.HTTPFirewallPolicy, error) {
	client, err := rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := client.HTTPFirewallPolicyRPC().FindEnabledFirewallPolicy(ctx, &pb.FindEnabledFirewallPolicyRequest{FirewallPolicyId: policyId})
	if err != nil {
		return nil, err
	}
	return resp.FirewallPolicy, nil
}

// 查找缓存策略配置
func (this *HTTPFirewallPolicyDAO) FindEnabledPolicyConfig(ctx context.Context, policyId int64) (*firewallconfigs.HTTPFirewallPolicy, error) {
	client, err := rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := client.HTTPFirewallPolicyRPC().FindEnabledFirewallPolicyConfig(ctx, &pb.FindEnabledFirewallPolicyConfigRequest{FirewallPolicyId: policyId})
	if err != nil {
		return nil, err
	}
	if len(resp.FirewallPolicyJSON) == 0 {
		return nil, nil
	}
	firewallPolicy := &firewallconfigs.HTTPFirewallPolicy{}
	err = json.Unmarshal(resp.FirewallPolicyJSON, firewallPolicy)
	if err != nil {
		return nil, err
	}
	return firewallPolicy, nil
}
