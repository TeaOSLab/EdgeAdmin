// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.
//go:build !plus

package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/maps"
)

type PlansAction struct {
	actionutils.ParentAction
}

func (this *PlansAction) RunPost(params struct {
	UserId   int64
	ServerId int64
}) {
	this.Data["plans"] = []maps.Map{}
	this.Success()
}
