// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.
//go:build !plus

package helpers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/iwind/TeaGo/maps"
)

func FindAllMenuMaps(langCode string, nodeLogsType string, countUnreadNodeLogs int64, countUnreadIPItems int64) []maps.Map {
	return []maps.Map{
		{
			"code":   "dashboard",
			"module": configloaders.AdminModuleCodeDashboard,
			"name":   langs.Message(langCode, codes.AdminMenuDashboard),
			"icon":   "dashboard",
		},
		{
			"code":     "servers",
			"module":   configloaders.AdminModuleCodeServer,
			"name":     langs.Message(langCode, codes.AdminMenuServers),
			"subtitle": "",
			"icon":     "clone outsize",
			"subItems": []maps.Map{
				{
					"name": langs.Message(langCode, codes.AdminMenuServerAccessLogs),
					"url":  "/servers/logs",
					"code": "log",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuServerCerts),
					"url":  "/servers/certs",
					"code": "cert",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuServerGroups),
					"url":  "/servers/groups",
					"code": "group",
				},
				{
					"name": "-",
					"url":  "",
					"code": "",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuServerCachePolicies),
					"url":  "/servers/components/cache",
					"code": "cache",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuServerPurgeFetchCaches),
					"url":  "/servers/components/cache/batch",
					"code": "cacheBatch",
				},
				{
					"name": "-",
					"url":  "",
					"code": "",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuServerWAFPolicies),
					"url":  "/servers/components/waf",
					"code": "waf",
				},
				{
					"name":  langs.Message(langCode, codes.AdminMenuServerIPLists),
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
					"name": langs.Message(langCode, codes.AdminMenuServerMetrics),
					"url":  "/servers/metrics",
					"code": "metric",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuServerGlobalSettings),
					"url":  "/servers/components",
					"code": "global",
				},
			},
		},
		{
			"code":     "clusters",
			"module":   configloaders.AdminModuleCodeNode,
			"name":     langs.Message(langCode, codes.AdminMenuNodes),
			"subtitle": "",
			"icon":     "cloud",
			"subItems": []maps.Map{
				{
					"name": langs.Message(langCode, codes.AdminMenuNodeClusters),
					"url":  "/clusters",
					"code": "cluster",
				},
				{
					"name":  langs.Message(langCode, codes.AdminMenuNodeLogs),
					"url":   "/clusters/logs?type=" + nodeLogsType,
					"code":  "log",
					"badge": countUnreadNodeLogs,
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuNodeRegions),
					"url":  "/clusters/regions",
					"code": "region",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuNodeSSHGrants),
					"url":  "/clusters/grants",
					"code": "grant",
				},
			},
		},
		{
			"code":     "dns",
			"module":   configloaders.AdminModuleCodeDNS,
			"name":     langs.Message(langCode, codes.AdminMenuDNS),
			"subtitle": "",
			"icon":     "globe",
			"subItems": []maps.Map{
				{
					"name": langs.Message(langCode, codes.AdminMenuDNSClusters),
					"url":  "/dns",
					"code": "cluster",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuDNSProviders),
					"url":  "/dns/providers",
					"code": "provider",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuDNSIssues),
					"url":  "/dns/issues",
					"code": "issue",
				},
			},
		},
		{
			"code":   "users",
			"module": configloaders.AdminModuleCodeUser,
			"name":   langs.Message(langCode, codes.AdminMenuUsers),
			"icon":   "users",
		},
		{
			"code":     "admins",
			"module":   configloaders.AdminModuleCodeAdmin,
			"name":     langs.Message(langCode, codes.AdminMenuAdmins),
			"subtitle": "",
			"icon":     "user secret",
		},
		{
			"code":   "log",
			"module": configloaders.AdminModuleCodeLog,
			"name":   langs.Message(langCode, codes.AdminMenuLogs),
			"icon":   "history",
		},
		{
			"code":     "settings",
			"module":   configloaders.AdminModuleCodeSetting,
			"name":     langs.Message(langCode, codes.AdminMenuSettings),
			"subtitle": "",
			"icon":     "setting",
			"subItems": []maps.Map{
				{
					"name": langs.Message(langCode, codes.AdminMenuSettingBasicSettings),
					"url":  "/settings",
					"code": "basic",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenuSettingAdvancedSettings),
					"url":  "/settings/advanced",
					"code": "advanced",
				},
			},
		},
	}
}
