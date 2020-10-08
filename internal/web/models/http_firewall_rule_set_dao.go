package models

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
)

var SharedHTTPFirewallRuleSetDAO = new(HTTPFirewallRuleSetDAO)

type HTTPFirewallRuleSetDAO struct {
}

// 查找规则集配置
func (this *HTTPFirewallRuleSetDAO) FindRuleSetConfig(ctx context.Context, setId int64) (*firewallconfigs.HTTPFirewallRuleSet, error) {
	client, err := rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := client.HTTPFirewallRuleSetRPC().FindHTTPFirewallRuleSetConfig(ctx, &pb.FindHTTPFirewallRuleSetConfigRequest{FirewallRuleSetId: setId})
	if err != nil {
		return nil, err
	}
	if len(resp.FirewallRuleSetJSON) == 0 {
		return nil, err
	}
	config := &firewallconfigs.HTTPFirewallRuleSet{}
	err = json.Unmarshal(resp.FirewallRuleSetJSON, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
