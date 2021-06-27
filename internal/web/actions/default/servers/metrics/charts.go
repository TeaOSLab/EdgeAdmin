// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package metrics

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type ChartsAction struct {
	actionutils.ParentAction
}

func (this *ChartsAction) Init() {
	this.Nav("", "", "chart")
}

func (this *ChartsAction) RunGet(params struct {
	ItemId int64
}) {
	err := InitItem(this.Parent(), params.ItemId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}
