// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/rands"
)

type RandomNameAction struct {
	actionutils.ParentAction
}

func (this *RandomNameAction) RunPost(params struct{}) {
	this.Data["name"] = "cluster" + rands.HexString(8)

	this.Success()
}
