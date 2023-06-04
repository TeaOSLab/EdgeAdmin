// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.
//go:build !plus

package helpers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/iwind/TeaGo/maps"
)

func FindAllMenuMaps(nodeLogsType string, countUnreadNodeLogs int64, countUnreadIPItems int64) []maps.Map {
	return []maps.Map{
		{
			"code":   "dashboard",
			"module": configloaders.AdminModuleCodeDashboard,
			"name":   "数据看板",
			"icon":   "dashboard",
		},
		{
			"code":     "servers",
			"module":   configloaders.AdminModuleCodeServer,
			"name":     "网站列表",
			"subtitle": "",
			"icon":     "clone outsize",
			"subItems": []maps.Map{
				{
					"name": "访问日志",
					"url":  "/servers/logs",
					"code": "log",
				},
				{
					"name": "证书管理",
					"url":  "/servers/certs",
					"code": "cert",
				},
				{
					"name": "网站分组",
					"url":  "/servers/groups",
					"code": "group",
				},
				{
					"name": "-",
					"url":  "",
					"code": "",
				},
				{
					"name": "缓存策略",
					"url":  "/servers/components/cache",
					"code": "cache",
				},
				{
					"name": "刷新预热",
					"url":  "/servers/components/cache/batch",
					"code": "cacheBatch",
				},
				{
					"name": "-",
					"url":  "",
					"code": "",
				},
				{
					"name": "WAF策略",
					"url":  "/servers/components/waf",
					"code": "waf",
				},
				{
					"name":  "IP名单",
					"url":   "/servers/iplists",
					"code":  "iplist",
					"badge": countUnreadIPItems,
				},
				{
					"name": "-",
					"url":  "",
					"code": "",
				},
				{
					"name": "统计指标",
					"url":  "/servers/metrics",
					"code": "metric",
				},
				{
					"name": "通用设置",
					"url":  "/servers/components",
					"code": "global",
				},
			},
		},
		{
			"code":     "clusters",
			"module":   configloaders.AdminModuleCodeNode,
			"name":     "边缘节点",
			"subtitle": "集群列表",
			"icon":     "cloud",
			"subItems": []maps.Map{
				{
					"name":  "运行日志",
					"url":   "/clusters/logs?type=" + nodeLogsType,
					"code":  "log",
					"badge": countUnreadNodeLogs,
				},
				{
					"name": "SSH认证",
					"url":  "/clusters/grants",
					"code": "grant",
				},
				{
					"name": "区域设置",
					"url":  "/clusters/regions",
					"code": "region",
				},
			},
		},
		{
			"code":     "dns",
			"module":   configloaders.AdminModuleCodeDNS,
			"name":     "域名解析",
			"subtitle": "集群列表",
			"icon":     "globe",
			"subItems": []maps.Map{
				{
					"name": "DNS服务商",
					"url":  "/dns/providers",
					"code": "provider",
				},
				{
					"name": "问题修复",
					"url":  "/dns/issues",
					"code": "issue",
				},
			},
		},
		{
			"code":   "users",
			"module": configloaders.AdminModuleCodeUser,
			"name":   "平台用户",
			"icon":   "users",
		},
		{
			"code":     "admins",
			"module":   configloaders.AdminModuleCodeAdmin,
			"name":     "系统用户",
			"subtitle": "用户列表",
			"icon":     "user secret",
		},
		{
			"code":   "log",
			"module": configloaders.AdminModuleCodeLog,
			"name":   "日志审计",
			"icon":   "history",
		},
		{
			"code":     "settings",
			"module":   configloaders.AdminModuleCodeSetting,
			"name":     "系统设置",
			"subtitle": "基本设置",
			"icon":     "setting",
			"subItems": []maps.Map{
				{
					"name": "高级设置",
					"url":  "/settings/advanced",
					"code": "advanced",
				},
			},
		},
	}
}
