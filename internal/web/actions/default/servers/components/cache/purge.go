package cache

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/cache/cacheutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net/http"
	"strconv"
	"strings"
)

type PurgeAction struct {
	actionutils.ParentAction
}

func (this *PurgeAction) Init() {
	this.Nav("", "", "purge")
}

func (this *PurgeAction) RunGet(params struct {
	CachePolicyId int64
}) {
	// 默认的集群ID
	cookie, err := this.Request.Cookie("cache_cluster_id")
	if cookie != nil && err == nil {
		this.Data["clusterId"] = types.Int64(cookie.Value)
	}

	// 集群列表
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClustersWithHTTPCachePolicyId(this.AdminContext(), &pb.FindAllEnabledNodeClustersWithHTTPCachePolicyIdRequest{HttpCachePolicyId: params.CachePolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range clustersResp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}

func (this *PurgeAction) RunPost(params struct {
	CachePolicyId int64
	ClusterId     int64
	KeyType       string
	Keys          string
	Must          *actions.Must
}) {
	// 创建日志
	defer this.CreateLogInfo(codes.ServerCachePolicy_LogPurgeCaches, params.CachePolicyId)

	// 记录clusterId
	this.AddCookie(&http.Cookie{
		Name:  "cache_cluster_id",
		Value: strconv.FormatInt(params.ClusterId, 10),
	})

	cachePolicyResp, err := this.RPC().HTTPCachePolicyRPC().FindEnabledHTTPCachePolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPCachePolicyConfigRequest{HttpCachePolicyId: params.CachePolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cachePolicyJSON := cachePolicyResp.HttpCachePolicyJSON
	if len(cachePolicyJSON) == 0 {
		this.Fail("找不到要操作的缓存策略")
	}

	if len(params.Keys) == 0 {
		this.Fail("请输入要删除的Key列表")
	}

	realKeys := []string{}
	for _, key := range strings.Split(params.Keys, "\n") {
		key = strings.TrimSpace(key)
		if len(key) == 0 {
			continue
		}
		if lists.ContainsString(realKeys, key) {
			continue
		}
		realKeys = append(realKeys, key)
	}
	// 校验Key
	validateResp, err := this.RPC().HTTPCacheTaskKeyRPC().ValidateHTTPCacheTaskKeys(this.AdminContext(), &pb.ValidateHTTPCacheTaskKeysRequest{Keys: realKeys})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var failKeyMaps = []maps.Map{}
	if len(validateResp.FailKeys) > 0 {
		for _, key := range validateResp.FailKeys {
			failKeyMaps = append(failKeyMaps, maps.Map{
				"key":    key.Key,
				"reason": cacheutils.KeyFailReason(key.ReasonCode),
			})
		}
	}
	this.Data["failKeys"] = failKeyMaps
	if len(failKeyMaps) > 0 {
		this.Fail("有" + types.String(len(failKeyMaps)) + "个Key无法完成操作，请删除后重试")
	}

	// 提交任务
	_, err = this.RPC().HTTPCacheTaskRPC().CreateHTTPCacheTask(this.AdminContext(), &pb.CreateHTTPCacheTaskRequest{
		Type:    "purge",
		KeyType: params.KeyType,
		Keys:    realKeys,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
