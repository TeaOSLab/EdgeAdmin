// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package boards

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
)

type WafAction struct {
	actionutils.ParentAction
}

func (this *WafAction) Init() {
	this.Nav("", "", "waf")
}

func (this *WafAction) RunGet(params struct{}) {
	if !teaconst.IsPlus {
		this.RedirectURL("/dashboard")
		return
	}

	this.Show()
}
