// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package cluster

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type UpgradeRemoteAction struct {
	actionutils.ParentAction
}

func (this *UpgradeRemoteAction) Init() {
	this.Nav("", "", "")
}

func (this *UpgradeRemoteAction) RunGet(params struct{}) {
	this.Show()
}
