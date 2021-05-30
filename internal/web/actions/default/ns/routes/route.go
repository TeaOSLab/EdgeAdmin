// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type RouteAction struct {
	actionutils.ParentAction
}

func (this *RouteAction) Init() {
	this.Nav("", "", "")
}

func (this *RouteAction) RunGet(params struct{}) {
	this.Show()
}
