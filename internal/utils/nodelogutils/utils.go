// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.
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
			"name": langs.Message(langCode, codes.Log_TagListener),
			"code": "LISTENER",
		},
		{
			"name": langs.Message(langCode, codes.Log_TagWAF),
			"code": "WAF",
		},
		{
			"name": langs.Message(langCode, codes.Log_TagAccessLog),
			"code": "ACCESS_LOG",
		},
	}
}
