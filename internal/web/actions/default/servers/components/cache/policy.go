package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/cache/cacheutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

type PolicyAction struct {
	actionutils.ParentAction
}

func (this *PolicyAction) Init() {
	this.Nav("", "", "index")
}

func (this *PolicyAction) RunGet(params struct {
	CachePolicyId int64
}) {
	cachePolicy, err := cacheutils.FindCachePolicy(this.Parent(), params.CachePolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["cachePolicy"] = cachePolicy

	this.Data["typeName"] = serverconfigs.FindCachePolicyStorageName(cachePolicy.Type)

	this.Show()
}
