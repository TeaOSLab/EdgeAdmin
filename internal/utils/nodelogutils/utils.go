// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.
//go:build community
// +build community

package nodelogutils

import (
	"github.com/iwind/TeaGo/maps"
)

// FindCommonTags 查找常用的标签
func FindNodeCommonTags() []maps.Map {
	return []maps.Map{
		{
			"name": "端口监听",
			"code": "LISTENER",
		},
		{
			"name": "WAF",
			"code": "WAF",
		},
		{
			"name": "访问日志",
			"code": "ACCESS_LOG",
		},
	}
}
