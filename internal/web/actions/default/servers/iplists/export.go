// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type ExportAction struct {
	actionutils.ParentAction
}

func (this *ExportAction) Init() {
	this.Nav("", "", "export")
}

func (this *ExportAction) RunGet(params struct {
	ListId int64
}) {
	err := InitIPList(this.Parent(), params.ListId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}
