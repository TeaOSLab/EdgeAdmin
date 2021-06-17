// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package access

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct{}) {
	this.Show()
}
