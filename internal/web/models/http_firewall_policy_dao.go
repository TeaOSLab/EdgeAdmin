package models

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/ipconfigs"
)

var SharedHTTPFirewallPolicyDAO = new(HTTPFirewallPolicyDAO)

// WAF策略相关
type HTTPFirewallPolicyDAO struct {
}

// 查找WAF策略基本信息
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

// 查找WAF策略配置
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

// 查找WAF的Inbound
func (this *HTTPFirewallPolicyDAO) FindEnabledPolicyInboundConfig(ctx context.Context, policyId int64) (*firewallconfigs.HTTPFirewallInboundConfig, error) {
	config, err := this.FindEnabledPolicyConfig(ctx, policyId)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("not found")
	}
	return config.Inbound, nil
}

// 根据类型查找WAF的IP名单
func (this *HTTPFirewallPolicyDAO) FindEnabledPolicyIPListIdWithType(ctx context.Context, policyId int64, listType ipconfigs.IPListType) (int64, error) {
	switch listType {
	case ipconfigs.IPListTypeWhite:
		return this.FindEnabledPolicyWhiteIPListId(ctx, policyId)
	case ipconfigs.IPListTypeBlack:
		return this.FindEnabledPolicyBlackIPListId(ctx, policyId)
	default:
		return 0, errors.New("invalid ip list type '" + listType + "'")
	}
}

// 查找WAF的白名单
func (this *HTTPFirewallPolicyDAO) FindEnabledPolicyWhiteIPListId(ctx context.Context, policyId int64) (int64, error) {
	client, err := rpc.SharedRPC()
	if err != nil {
		return 0, err
	}

	config, err := this.FindEnabledPolicyConfig(ctx, policyId)
	if err != nil {
		return 0, err
	}
	if config == nil {
		return 0, errors.New("not found")
	}
	if config.Inbound == nil {
		config.Inbound = &firewallconfigs.HTTPFirewallInboundConfig{IsOn: true}
	}
	if config.Inbound.WhiteListRef == nil || config.Inbound.WhiteListRef.ListId == 0 {
		createResp, err := client.IPListRPC().CreateIPList(ctx, &pb.CreateIPListRequest{
			Type:        "white",
			Name:        "白名单",
			Code:        "white",
			TimeoutJSON: nil,
		})
		if err != nil {
			return 0, err
		}
		listId := createResp.IpListId
		config.Inbound.WhiteListRef = &ipconfigs.IPListRef{
			IsOn:   true,
			ListId: listId,
		}
		inboundJSON, err := json.Marshal(config.Inbound)
		if err != nil {
			return 0, err
		}
		_, err = client.HTTPFirewallPolicyRPC().UpdateHTTPFirewallInboundConfig(ctx, &pb.UpdateHTTPFirewallInboundConfigRequest{
			FirewallPolicyId: policyId,
			InboundJSON:      inboundJSON,
		})
		if err != nil {
			return 0, err
		}
		return listId, nil
	}

	return config.Inbound.WhiteListRef.ListId, nil
}

// 查找WAF的黑名单
func (this *HTTPFirewallPolicyDAO) FindEnabledPolicyBlackIPListId(ctx context.Context, policyId int64) (int64, error) {
	client, err := rpc.SharedRPC()
	if err != nil {
		return 0, err
	}

	config, err := this.FindEnabledPolicyConfig(ctx, policyId)
	if err != nil {
		return 0, err
	}
	if config == nil {
		return 0, errors.New("not found")
	}
	if config.Inbound == nil {
		config.Inbound = &firewallconfigs.HTTPFirewallInboundConfig{IsOn: true}
	}
	if config.Inbound.BlackListRef == nil || config.Inbound.BlackListRef.ListId == 0 {
		createResp, err := client.IPListRPC().CreateIPList(ctx, &pb.CreateIPListRequest{
			Type:        "black",
			Name:        "黑名单",
			Code:        "black",
			TimeoutJSON: nil,
		})
		if err != nil {
			return 0, err
		}
		listId := createResp.IpListId
		config.Inbound.BlackListRef = &ipconfigs.IPListRef{
			IsOn:   true,
			ListId: listId,
		}
		inboundJSON, err := json.Marshal(config.Inbound)
		if err != nil {
			return 0, err
		}
		_, err = client.HTTPFirewallPolicyRPC().UpdateHTTPFirewallInboundConfig(ctx, &pb.UpdateHTTPFirewallInboundConfigRequest{
			FirewallPolicyId: policyId,
			InboundJSON:      inboundJSON,
		})
		if err != nil {
			return 0, err
		}
		return listId, nil
	}

	return config.Inbound.BlackListRef.ListId, nil
}
