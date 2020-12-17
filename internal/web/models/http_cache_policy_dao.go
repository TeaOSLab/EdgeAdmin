package models

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

var SharedHTTPCachePolicyDAO = new(HTTPCachePolicyDAO)

type HTTPCachePolicyDAO struct {
	BaseDAO
}

// 查找缓存策略配置
func (this *HTTPCachePolicyDAO) FindEnabledHTTPCachePolicyConfig(ctx context.Context, cachePolicyId int64) (*serverconfigs.HTTPCachePolicy, error) {
	resp, err := this.RPC().HTTPCachePolicyRPC().FindEnabledHTTPCachePolicyConfig(ctx, &pb.FindEnabledHTTPCachePolicyConfigRequest{HttpCachePolicyId: cachePolicyId})
	if err != nil {
		return nil, err
	}
	if len(resp.HttpCachePolicyJSON) == 0 {
		return nil, nil
	}
	config := &serverconfigs.HTTPCachePolicy{}
	err = json.Unmarshal(resp.HttpCachePolicyJSON, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// 查找缓存策略信息
func (this *HTTPCachePolicyDAO) FindEnabledHTTPCachePolicy(ctx context.Context, cachePolicyId int64) (*pb.HTTPCachePolicy, error) {
	resp, err := this.RPC().HTTPCachePolicyRPC().FindEnabledHTTPCachePolicy(ctx, &pb.FindEnabledHTTPCachePolicyRequest{
		HttpCachePolicyId: cachePolicyId,
	})
	if err != nil {
		return nil, err
	}
	return resp.HttpCachePolicy, nil
}

// 根据服务ID查找缓存策略
func (this *HTTPCachePolicyDAO) FindEnabledHTTPCachePolicyWithServerId(ctx context.Context, serverId int64) (*pb.HTTPCachePolicy, error) {
	serverResp, err := this.RPC().ServerRPC().FindEnabledServer(ctx, &pb.FindEnabledServerRequest{ServerId: serverId})
	if err != nil {
		return nil, err
	}
	server := serverResp.Server
	if server == nil {
		return nil, nil
	}
	if server.NodeCluster == nil {
		return nil, nil
	}
	clusterId := server.NodeCluster.Id
	cluster, err := SharedNodeClusterDAO.FindEnabledNodeCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, nil
	}
	if cluster.HttpCachePolicyId == 0 {
		return nil, nil
	}
	return SharedHTTPCachePolicyDAO.FindEnabledHTTPCachePolicy(ctx, cluster.HttpCachePolicyId)
}
