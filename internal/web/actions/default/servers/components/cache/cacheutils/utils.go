package cacheutils

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

// FindCachePolicyNameWithoutError 查找缓存策略名称并忽略错误
func FindCachePolicyNameWithoutError(parent *actionutils.ParentAction, cachePolicyId int64) string {
	policy, err := FindCachePolicy(parent, cachePolicyId)
	if err != nil {
		return ""
	}
	if policy == nil {
		return ""
	}
	return policy.Name
}

// FindCachePolicy 查找缓存策略配置
func FindCachePolicy(parent *actionutils.ParentAction, cachePolicyId int64) (*serverconfigs.HTTPCachePolicy, error) {
	resp, err := parent.RPC().HTTPCachePolicyRPC().FindEnabledHTTPCachePolicyConfig(parent.AdminContext(), &pb.FindEnabledHTTPCachePolicyConfigRequest{HttpCachePolicyId: cachePolicyId})
	if err != nil {
		return nil, err
	}
	if len(resp.HttpCachePolicyJSON) == 0 {
		return nil, errors.New("cache policy not found")
	}
	var config = &serverconfigs.HTTPCachePolicy{}
	err = json.Unmarshal(resp.HttpCachePolicyJSON, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// KeyFailReason Key相关失败原因
func KeyFailReason(reasonCode string) string {
	switch reasonCode {
	case "requireKey":
		return "空的Key"
	case "requireDomain":
		return "找不到Key对应的域名"
	case "requireServer":
		return "找不到Key对应的网站"
	case "requireUser":
		return "该域名不属于当前用户"
	case "requireClusterId":
		return "该网站没有部署到集群"
	}
	return "未知错误"
}
