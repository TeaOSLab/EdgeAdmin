package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	IsReverse bool
}) {
	this.Data["isReverse"] = params.IsReverse

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	CacheRefJSON []byte

	CondType              string
	CondJSON              []byte
	CondIsCaseInsensitive bool

	Must *actions.Must
}) {
	var cacheRef = &serverconfigs.HTTPCacheRef{}
	err := json.Unmarshal(params.CacheRefJSON, cacheRef)
	if err != nil {
		this.Fail("解析条件出错：" + err.Error() + ", JSON: " + string(params.CacheRefJSON))
		return
	}

	if len(params.CondJSON) > 0 {
		var cond = &shared.HTTPRequestCond{}
		err = json.Unmarshal(params.CondJSON, cond)
		if err != nil {
			this.Fail("解析条件出错：" + err.Error() + ", JSON: " + string(params.CondJSON))
			return
		}
		cond.Type = params.CondType
		cond.IsCaseInsensitive = params.CondIsCaseInsensitive
		cacheRef.SimpleCond = cond

		// 将组合条件置为空
		cacheRef.Conds = &shared.HTTPRequestCondsConfig{}
	}

	err = cacheRef.Init()
	if err != nil {
		this.Fail("解析条件出错：" + err.Error())
		return
	}

	if len(cacheRef.Key) == 0 {
		this.Fail("请输入缓存Key")
	}

	if (cacheRef.Conds == nil || len(cacheRef.Conds.Groups) == 0) && cacheRef.SimpleCond == nil {
		this.Fail("请填写匹配条件分组")
	}

	this.Data["cacheRef"] = cacheRef

	cacheRefClone, err := utils.JSONClone(cacheRef)
	if err != nil {
		this.Fail(err.Error())
	}
	err = cacheRefClone.(*serverconfigs.HTTPCacheRef).Init()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["cacheRef"] = cacheRef

	this.Success()
}
