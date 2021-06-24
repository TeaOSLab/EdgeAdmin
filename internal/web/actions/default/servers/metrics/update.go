// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package metrics

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateAction) RunGet(params struct{}) {
	this.Show()
}
