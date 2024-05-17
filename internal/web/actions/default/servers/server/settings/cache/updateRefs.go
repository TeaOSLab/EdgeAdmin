// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package cache

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

type UpdateRefsAction struct {
	actionutils.ParentAction
}

func (this *UpdateRefsAction) RunPost(params struct {
	WebId    int64
	RefsJSON []byte
}) {
	// 日志
	defer this.CreateLogInfo(codes.ServerCache_LogUpdateCacheSettings, params.WebId)

	this.Data["isUpdated"] = false

	webConfigResp, err := this.RPC().HTTPWebRPC().FindEnabledHTTPWebConfig(this.AdminContext(), &pb.FindEnabledHTTPWebConfigRequest{HttpWebId: params.WebId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var webConfig = &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(webConfigResp.HttpWebJSON, webConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["cacheConfig"] = webConfig.Cache

	// 校验配置
	var cacheConfig = webConfig.Cache
	if cacheConfig == nil {
		this.Success()
		return
	}

	var refs = []*serverconfigs.HTTPCacheRef{}
	err = json.Unmarshal(params.RefsJSON, &refs)
	if err != nil {
		this.ErrorPage(errors.New("decode refs json failed: " + err.Error()))
		return
	}
	cacheConfig.CacheRefs = refs

	err = cacheConfig.Init()
	if err != nil {
		this.Fail("检查配置失败：" + err.Error())
	}

	// 去除不必要的部分
	for _, cacheRef := range cacheConfig.CacheRefs {
		cacheRef.CachePolicy = nil
	}

	cacheJSON, err := json.Marshal(cacheConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebCache(this.AdminContext(), &pb.UpdateHTTPWebCacheRequest{
		HttpWebId: params.WebId,
		CacheJSON: cacheJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["isUpdated"] = true
	this.Success()
}
