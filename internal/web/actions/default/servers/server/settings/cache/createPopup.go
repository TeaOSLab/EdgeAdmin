package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
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

	Must *actions.Must
}) {
	cacheRef := &serverconfigs.HTTPCacheRef{}
	err := json.Unmarshal(params.CacheRefJSON, cacheRef)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if len(cacheRef.Key) == 0 {
		this.Fail("请输入缓存Key")
	}

	if cacheRef.Conds == nil || len(cacheRef.Conds.Groups) == 0 {
		this.Fail("请填写匹配条件分组")
	}

	err = cacheRef.Init()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["cacheRef"] = cacheRef

	this.Success()
}
