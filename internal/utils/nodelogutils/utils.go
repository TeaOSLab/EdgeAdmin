// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.
//go:build !plus

package nodelogutils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/langs"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/iwind/TeaGo/maps"
)

// FindCommonTags 查找常用的标签
func FindNodeCommonTags(langCode langs.LangCode) []maps.Map {
	return []maps.Map{
		{
			"name": langs.Message(langCode, codes.AdminLogTagListener),
			"code": "LISTENER",
		},
		{
			"name": langs.Message(langCode, codes.AdminLogTagWAF),
			"code": "WAF",
		},
		{
			"name": langs.Message(langCode, codes.AdminLogTagAccessLog),
			"code": "ACCESS_LOG",
		},
	}
}
