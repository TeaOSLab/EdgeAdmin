// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package setup

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

var currentStatusText = ""

type StatusAction struct {
	actionutils.ParentAction
}

func (this *StatusAction) RunPost(params struct{}) {
	this.Data["statusText"] = currentStatusText
	this.Success()
}
