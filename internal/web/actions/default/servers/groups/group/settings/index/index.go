// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package reverseProxy

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/types"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
}) {
	if teaconst.IsPlus {
		this.RedirectURL("/servers/groups/group/settings/web?groupId=" + types.String(params.GroupId))
	} else {
		this.RedirectURL("/servers/groups/group/settings/httpReverseProxy?groupId=" + types.String(params.GroupId))
	}
}
