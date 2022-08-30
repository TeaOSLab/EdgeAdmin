// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package access

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/rands"
)

type RandomAction struct {
	actionutils.ParentAction
}

func (this *RandomAction) RunPost(params struct{}) {
	this.Data["random"] = rands.HexString(32)

	this.Success()
}
