package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("cache")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	webConfigResp, err := this.RPC().ServerRPC().FindAndInitServerWebConfig(this.AdminContext(), &pb.FindAndInitServerWebRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	webConfig := &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(webConfigResp.Config, webConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["cacheConfig"] = webConfig.CacheRef

	// 所有缓存策略
	cachePoliciesResp, err := this.RPC().HTTPCachePolicyRPC().FindAllEnabledHTTPCachePolicies(this.AdminContext(), &pb.FindAllEnabledHTTPCachePoliciesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyMaps := []maps.Map{}
	for _, policy := range cachePoliciesResp.CachePolicies {
		policyMaps = append(policyMaps, maps.Map{
			"id":   policy.Id,
			"name": policy.Name,
			"isOn": policy.IsOn,
		})
	}
	this.Data["policies"] = policyMaps

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId     int64
	CacheJSON []byte

	Must *actions.Must
}) {
	// TODO 校验配置

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPCache(this.AdminContext(), &pb.UpdateHTTPCacheRequest{
		WebId:     params.WebId,
		CacheJSON: params.CacheJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
