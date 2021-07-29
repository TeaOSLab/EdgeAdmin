// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package policyutils

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

// InitPolicy 初始化访问日志策略
func InitPolicy(parent *actionutils.ParentAction, policyId int64) error {
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	policyResp, err := rpcClient.HTTPAccessLogPolicyRPC().FindEnabledHTTPAccessLogPolicy(parent.AdminContext(), &pb.FindEnabledHTTPAccessLogPolicyRequest{HttpAccessLogPolicyId: policyId})
	if err != nil {
		return err
	}
	var policy = policyResp.HttpAccessLogPolicy
	if policy == nil {
		return errors.New("can not find policy '" + types.String(policyId) + "'")
	}

	var options = maps.Map{}
	if len(policy.OptionsJSON) > 0 {
		err = json.Unmarshal(policy.OptionsJSON, &options)
		if err != nil {
			return err
		}
	}

	parent.Data["policy"] = maps.Map{
		"id":       policy.Id,
		"name":     policy.Name,
		"type":     policy.Type,
		"typeName": serverconfigs.FindAccessLogStorageTypeName(policy.Type),
		"isOn":     policy.IsOn,
		"isPublic": policy.IsPublic,
		"options":  options,
	}
	return nil
}
