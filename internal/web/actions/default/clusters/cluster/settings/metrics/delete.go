// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package metrics

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct{}) {
	this.Success()
}
