// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package cluster

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"strconv"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	if teaconst.IsPlus {
		this.RedirectURL("/clusters/cluster/boards?clusterId=" + strconv.FormatInt(params.ClusterId, 10))
	} else {
		this.RedirectURL("/clusters/cluster/nodes?clusterId=" + strconv.FormatInt(params.ClusterId, 10))
	}
}
