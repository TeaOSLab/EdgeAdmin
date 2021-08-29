// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/messageconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"strings"
)

type PreheatAction struct {
	actionutils.ParentAction
}

func (this *PreheatAction) Init() {
	this.Nav("", "setting", "preheat")
	this.SecondMenu("cache")
}

func (this *PreheatAction) RunGet(params struct {
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

func (this *PreheatAction) RunPost(params struct {
	ServerId int64
	WebId    int64
	Keys     string

	Must *actions.Must
}) {

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "预热服务 %d 缓存", params.ServerId)

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

	// 发送命令
	msg := &messageconfigs.PreheatCacheMessage{
		CachePolicyJSON: cachePolicyJSON,
		Keys:            realKeys,
	}
	results, err := nodeutils.SendMessageToCluster(this.AdminContext(), clusterId, messageconfigs.MessageCodePreheatCache, msg, 300)
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
