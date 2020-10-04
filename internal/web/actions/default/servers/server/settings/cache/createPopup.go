package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/models"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	// 缓存策略列表
	cachePoliciesResp, err := this.RPC().HTTPCachePolicyRPC().FindAllEnabledHTTPCachePolicies(this.AdminContext(), &pb.FindAllEnabledHTTPCachePoliciesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cachePolicyMaps := []maps.Map{}
	for _, cachePolicy := range cachePoliciesResp.CachePolicies {
		cachePolicyMaps = append(cachePolicyMaps, maps.Map{
			"id":   cachePolicy.Id,
			"name": cachePolicy.Name,
		})
	}
	this.Data["cachePolicies"] = cachePolicyMaps

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	CachePolicyId int64
	CacheRefJSON  []byte

	Must *actions.Must
}) {
	if params.CachePolicyId <= 0 {
		this.Fail("请选择要使用的缓存策略")
	}

	cachePolicy, err := models.SharedHTTPCachePolicyDAO.FindEnabledCachePolicyConfig(this.AdminContext(), params.CachePolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if cachePolicy == nil {
		this.Fail("找不到你要使用的缓存策略")
	}

	cacheRef := &serverconfigs.HTTPCacheRef{}
	err = json.Unmarshal(params.CacheRefJSON, cacheRef)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if len(cacheRef.Key) == 0 {
		this.Fail("请输入缓存Key")
	}

	cacheRef.CachePolicyId = cachePolicy.Id
	cacheRef.CachePolicy = cachePolicy

	err = cacheRef.Init()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["cacheRef"] = cacheRef

	this.Success()
}
