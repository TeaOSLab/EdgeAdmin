// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package lang

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"net/http"
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

	this.AddCookie(&http.Cookie{
		Name:  "edgelang",
		Value: langCode,
		Path:  "/",
	})

	this.Success()
}
