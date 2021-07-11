// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package boards

import "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"

type DnsAction struct {
	actionutils.ParentAction
}

func (this *DnsAction) Init() {
	this.Nav("", "", "dns")
}

func (this *DnsAction) RunGet(params struct{}) {
	this.Show()
}
