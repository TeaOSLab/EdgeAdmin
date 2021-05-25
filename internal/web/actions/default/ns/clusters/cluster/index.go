// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package cluster

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	this.Show()
}
