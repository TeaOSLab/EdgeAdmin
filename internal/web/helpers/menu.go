// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.
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
			"name":   langs.Message(langCode, codes.AdminMenu_Dashboard),
			"icon":   "dashboard",
		},
		{
			"code":     "servers",
			"module":   configloaders.AdminModuleCodeServer,
			"name":     langs.Message(langCode, codes.AdminMenu_Servers),
			"subtitle": "",
			"icon":     "clone outsize",
			"subItems": []maps.Map{
				{
					"name": langs.Message(langCode, codes.AdminMenu_ServerAccessLogs),
					"url":  "/servers/logs",
					"code": "log",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenu_ServerCerts),
					"url":  "/servers/certs",
					"code": "cert",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenu_ServerGroups),
					"url":  "/servers/groups",
					"code": "group",
				},
				{
					"name": "-",
					"url":  "",
					"code": "",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenu_ServerCachePolicies),
					"url":  "/servers/components/cache",
					"code": "cache",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenu_ServerPurgeFetchCaches),
					"url":  "/servers/components/cache/batch",
					"code": "cacheBatch",
				},
				{
					"name": "-",
					"url":  "",
					"code": "",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenu_ServerWAFPolicies),
					"url":  "/servers/components/waf",
					"code": "waf",
				},
				{
					"name":  langs.Message(langCode, codes.AdminMenu_ServerIPLists),
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
					"name": langs.Message(langCode, codes.AdminMenu_ServerMetrics),
					"url":  "/servers/metrics",
					"code": "metric",
				},
			},
		},
		{
			"code":     "clusters",
			"module":   configloaders.AdminModuleCodeNode,
			"name":     langs.Message(langCode, codes.AdminMenu_Nodes),
			"subtitle": "",
			"icon":     "cloud",
			"subItems": []maps.Map{
				{
					"name": langs.Message(langCode, codes.AdminMenu_NodeClusters),
					"url":  "/clusters",
					"code": "cluster",
				},
				{
					"name":  langs.Message(langCode, codes.AdminMenu_NodeLogs),
					"url":   "/clusters/logs?type=" + nodeLogsType,
					"code":  "log",
					"badge": countUnreadNodeLogs,
				},
				{
					"name": langs.Message(langCode, codes.AdminMenu_NodeRegions),
					"url":  "/clusters/regions",
					"code": "region",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenu_NodeSSHGrants),
					"url":  "/clusters/grants",
					"code": "grant",
				},
			},
		},
		{
			"code":     "dns",
			"module":   configloaders.AdminModuleCodeDNS,
			"name":     langs.Message(langCode, codes.AdminMenu_DNS),
			"subtitle": "",
			"icon":     "globe",
			"subItems": []maps.Map{
				{
					"name": langs.Message(langCode, codes.AdminMenu_DNSClusters),
					"url":  "/dns",
					"code": "cluster",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenu_DNSProviders),
					"url":  "/dns/providers",
					"code": "provider",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenu_DNSIssues),
					"url":  "/dns/issues",
					"code": "issue",
				},
			},
		},
		{
			"code":   "users",
			"module": configloaders.AdminModuleCodeUser,
			"name":   langs.Message(langCode, codes.AdminMenu_Users),
			"icon":   "users",
			"subItems": []maps.Map{
				{
					"name": langs.Message(langCode, codes.AdminMenu_UserList),
					"url":  "/users",
					"code": "users",
				},
			},
		},
		{
			"code":     "admins",
			"module":   configloaders.AdminModuleCodeAdmin,
			"name":     langs.Message(langCode, codes.AdminMenu_Admins),
			"subtitle": "",
			"icon":     "user secret",
		},
		{
			"code":   "log",
			"module": configloaders.AdminModuleCodeLog,
			"name":   langs.Message(langCode, codes.AdminMenu_Logs),
			"icon":   "history",
		},
		{
			"code":     "settings",
			"module":   configloaders.AdminModuleCodeSetting,
			"name":     langs.Message(langCode, codes.AdminMenu_Settings),
			"subtitle": "",
			"icon":     "setting",
			"subItems": []maps.Map{
				{
					"name": langs.Message(langCode, codes.AdminMenu_SettingBasicSettings),
					"url":  "/settings",
					"code": "basic",
				},
				{
					"name": langs.Message(langCode, codes.AdminMenu_SettingAdvancedSettings),
					"url":  "/settings/advanced",
					"code": "advanced",
				},
			},
		},
	}
}
