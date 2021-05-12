package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
}

func (this *IndexAction) RunGet(params struct {
	ServerId   int64
	LocationId int64
}) {
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithLocationId(this.AdminContext(), params.LocationId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["cacheConfig"] = webConfig.Cache

	// 当前集群的缓存策略
	cachePolicy, err := dao.SharedHTTPCachePolicyDAO.FindEnabledHTTPCachePolicyWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if cachePolicy != nil {
		this.Data["cachePolicy"] = maps.Map{
			"id":   cachePolicy.Id,
			"name": cachePolicy.Name,
			"isOn": cachePolicy.IsOn,
		}
	} else {
		this.Data["cachePolicy"] = nil
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId     int64
	CacheJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo("修改Web %d 的缓存设置", params.WebId)

	// 校验配置
	cacheConfig := &serverconfigs.HTTPCacheConfig{}
	err := json.Unmarshal(params.CacheJSON, cacheConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
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
		WebId:     params.WebId,
		CacheJSON: cacheJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
