// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package origins

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/ossconfigs"
	"github.com/iwind/TeaGo/maps"
)

func (this *AddPopupAction) getOSSHook() {
	this.Data["ossTypes"] = []maps.Map{}
	this.Data["ossBucketParams"] = []maps.Map{}
	this.Data["ossForm"] = ""
}

func (this *AddPopupAction) postOSSHook(protocol string) (config *ossconfigs.OSSConfig, goNext bool, err error) {
	goNext = true
	return
}
