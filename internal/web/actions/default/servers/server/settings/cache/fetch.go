// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/cache/cacheutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

type FetchAction struct {
	actionutils.ParentAction
}

func (this *FetchAction) Init() {
	this.Nav("", "setting", "fetch")
	this.SecondMenu("cache")
}

func (this *FetchAction) RunGet(params struct {
	ServerId int64
}) {
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["webConfig"] = webConfig

	this.Show()
}

func (this *FetchAction) RunPost(params struct {
	ServerId int64
	WebId    int64
	Keys     string

	Must *actions.Must
}) {

	// 创建日志
	defer this.CreateLogInfo(codes.ServerCache_LogFetchCaches, params.ServerId)

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithId(this.AdminContext(), params.WebId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if webConfig == nil {
		this.NotFound("httpWeb", params.WebId)
		return
	}
	var cache = webConfig.Cache
	if cache == nil || !cache.IsOn {
		this.Fail("当前没有开启缓存")
	}

	serverResp, err := this.RPC().ServerRPC().FindEnabledUserServerBasic(this.AdminContext(), &pb.FindEnabledUserServerBasicRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var server = serverResp.Server
	if server == nil || server.NodeCluster == nil {
		this.NotFound("server", params.ServerId)
		return
	}

	var clusterId = server.NodeCluster.Id

	clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{NodeClusterId: clusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var cluster = clusterResp.NodeCluster
	if cluster == nil {
		this.NotFound("nodeCluster", clusterId)
		return
	}
	var cachePolicyId = cluster.HttpCachePolicyId
	if cachePolicyId == 0 {
		this.Fail("当前集群没有设置缓存策略")
	}

	cachePolicyResp, err := this.RPC().HTTPCachePolicyRPC().FindEnabledHTTPCachePolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPCachePolicyConfigRequest{HttpCachePolicyId: cachePolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cachePolicyJSON := cachePolicyResp.HttpCachePolicyJSON
	if len(cachePolicyJSON) == 0 {
		this.Fail("找不到要操作的缓存策略")
	}

	if len(params.Keys) == 0 {
		this.Fail("请输入要预热的Key列表")
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
		Type:    "fetch",
		KeyType: "key",
		Keys:    realKeys,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
