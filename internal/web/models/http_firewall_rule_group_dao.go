package models

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
)

var SharedHTTPFirewallRuleGroupDAO = new(HTTPFirewallRuleGroupDAO)

type HTTPFirewallRuleGroupDAO struct {
}

// 查找分组配置
func (this *HTTPFirewallRuleGroupDAO) FindRuleGroupConfig(ctx context.Context, groupId int64) (*firewallconfigs.HTTPFirewallRuleGroup, error) {
	client, err := rpc.SharedRPC()
	if err != nil {
		return nil, err
	}

	groupResp, err := client.HTTPFirewallRuleGroupRPC().FindEnabledHTTPFirewallRuleGroupConfig(ctx, &pb.FindEnabledHTTPFirewallRuleGroupConfigRequest{FirewallRuleGroupId: groupId})
	if err != nil {
		return nil, err
	}

	if len(groupResp.FirewallRuleGroupJSON) == 0 {
		return nil, nil
	}

	groupConfig := &firewallconfigs.HTTPFirewallRuleGroup{}
	err = json.Unmarshal(groupResp.FirewallRuleGroupJSON, groupConfig)
	if err != nil {
		return nil, err
	}

	return groupConfig, nil
}
