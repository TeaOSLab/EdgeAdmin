// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package helpers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs"
	"github.com/iwind/TeaGo/actions"
)

type LangHelper struct {
}

func (this *LangHelper) Lang(actionPtr actions.ActionWrapper, messageCode langs.MessageCode, args ...any) string {
	var langCode = configloaders.FindAdminLang(actionPtr.Object().Session().GetInt64(teaconst.SessionAdminId))
	if len(langCode) == 0 {
		langCode = langs.ParseLangFromAction(actionPtr)
	}
	return langs.Message(langCode, messageCode, args...)
}
