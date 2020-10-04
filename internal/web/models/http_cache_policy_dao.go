package models

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

var SharedHTTPCachePolicyDAO = new(HTTPCachePolicyDAO)

type HTTPCachePolicyDAO struct {
}

// 查找缓存策略配置
func (this *HTTPCachePolicyDAO) FindEnabledCachePolicyConfig(ctx context.Context, cachePolicyId int64) (*serverconfigs.HTTPCachePolicy, error) {
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := rpcClient.HTTPCachePolicyRPC().FindEnabledHTTPCachePolicyConfig(ctx, &pb.FindEnabledHTTPCachePolicyConfigRequest{CachePolicyId: cachePolicyId})
	if err != nil {
		return nil, err
	}
	if len(resp.CachePolicyJSON) == 0 {
		return nil, nil
	}
	config := &serverconfigs.HTTPCachePolicy{}
	err = json.Unmarshal(resp.CachePolicyJSON, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
