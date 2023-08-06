package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/cache/cacheutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

type PolicyAction struct {
	actionutils.ParentAction
}

func (this *PolicyAction) Init() {
	this.Nav("", "", "index")
}

func (this *PolicyAction) RunGet(params struct {
	CachePolicyId int64
}) {
	cachePolicy, err := cacheutils.FindCachePolicy(this.Parent(), params.CachePolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["cachePolicy"] = cachePolicy

	// 预热超时时间
	this.Data["fetchTimeoutString"] = ""
	if cachePolicy.FetchTimeout != nil && cachePolicy.FetchTimeout.Count > 0 {
		this.Data["fetchTimeoutString"] = cachePolicy.FetchTimeout.Description()
	}

	this.Data["typeName"] = serverconfigs.FindCachePolicyStorageName(cachePolicy.Type)

	// 正在使用此策略的集群
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClustersWithHTTPCachePolicyId(this.AdminContext(), &pb.FindAllEnabledNodeClustersWithHTTPCachePolicyIdRequest{HttpCachePolicyId: params.CachePolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var clusterMaps = []maps.Map{}
	for _, cluster := range clustersResp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
