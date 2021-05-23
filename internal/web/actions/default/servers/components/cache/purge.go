package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/messageconfigs"
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
	Type          string
	Keys          string
	Must          *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "删除缓存，缓存策略：%d", params.CachePolicyId)

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

	// 发送命令
	msg := &messageconfigs.PurgeCacheMessage{
		CachePolicyJSON: cachePolicyJSON,
		Keys:            realKeys,
	}
	if params.Type == "prefix" {
		msg.Type = messageconfigs.PurgeCacheMessageTypeDir
	} else {
		msg.Type = messageconfigs.PurgeCacheMessageTypeFile
	}
	results, err := nodeutils.SendMessageToCluster(this.AdminContext(), params.ClusterId, messageconfigs.MessageCodePurgeCache, msg, 10)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	isAllOk := true
	for _, result := range results {
		if !result.IsOK {
			isAllOk = false
			break
		}
	}

	this.Data["isAllOk"] = isAllOk
	this.Data["results"] = results

	this.Success()
}
