// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package lang

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type SwitchAction struct {
	actionutils.ParentAction
}

func (this *SwitchAction) Init() {
	this.Nav("", "", "")
}

func (this *SwitchAction) RunPost(params struct{}) {
	var langCode = this.LangCode()
	if len(langCode) == 0 || langCode == "zh-cn" {
		langCode = "en-us"
	} else {
		langCode = "zh-cn"
	}

	configloaders.UpdateAdminLang(this.AdminId(), langCode)

	_, err := this.RPC().AdminRPC().UpdateAdminLang(this.AdminContext(), &pb.UpdateAdminLangRequest{LangCode: langCode})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
