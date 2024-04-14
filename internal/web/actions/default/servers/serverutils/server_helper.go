package serverutils

import (
	"encoding/json"
	"errors"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net/http"
	"strconv"
)

type ServerHelper struct {
	helpers.LangHelper
}

func NewServerHelper() *ServerHelper {
	return &ServerHelper{}
}

func (this *ServerHelper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "servers"

	// 左侧菜单
	this.createLeftMenu(action)
}

func (this *ServerHelper) createLeftMenu(action *actions.ActionObject) {
	// 初始化
	if !action.Data.Has("leftMenuItemIsDisabled") {
		action.Data["leftMenuItemIsDisabled"] = false
	}
	action.Data["leftMenuItems"] = []maps.Map{}
	var mainTab = action.Data["mainTab"]
	var secondMenuItem = action.Data["secondMenuItem"]

	serverId := action.ParamInt64("serverId")
	if serverId == 0 {
		return
	}
	var serverIdString = strconv.FormatInt(serverId, 10)
	action.Data["serverId"] = serverId

	// 读取server信息
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		logs.Error(err)
		return
	}

	serverResp, err := rpcClient.ServerRPC().FindEnabledServer(rpcClient.Context(action.Context.GetInt64(teaconst.SessionAdminId)), &pb.FindEnabledServerRequest{
		ServerId:       serverId,
		IgnoreSSLCerts: true,
	})
	if err != nil {
		logs.Error(err)
		return
	}
	var server = serverResp.Server
	if server == nil {
		logs.Error(errors.New("can not find the server"))
		return
	}

	// 初始化数据
	if !action.Data.Has("server") {
		if server.NodeCluster == nil {
			server.NodeCluster = &pb.NodeCluster{Id: 0}
		}
		action.Data["server"] = maps.Map{
			"id":        server.Id,
			"name":      server.Name,
			"clusterId": server.NodeCluster.Id,
		}
	}

	// 服务管理
	var serverConfig = &serverconfigs.ServerConfig{}
	err = json.Unmarshal(server.Config, serverConfig)
	if err != nil {
		logs.Error(err)
		return
	}

	// 协议簇
	var family = ""
	if serverConfig.IsHTTPFamily() {
		family = "http"
	} else if serverConfig.IsTCPFamily() {
		family = "tcp"
	} else if serverConfig.IsUDPFamily() {
		family = "udp"
	}
	action.Data["serverFamily"] = family

	// TABBAR
	var selectedTabbar = action.Data["mainTab"]
	var tabbar = actionutils.NewTabbar()
	tabbar.Add("", "", "/servers", "left arrow", false)
	if len(serverConfig.Name) > 0 {
		var item = tabbar.Add(serverConfig.Name, "", "/servers/server?serverId="+serverIdString, "angle right", true)
		item.IsTitle = true
	}

	if teaconst.IsPlus {
		tabbar.Add(this.Lang(action, codes.Server_TabDashboard), "", "/servers/server/boards?serverId="+serverIdString, "dashboard", selectedTabbar == "board")
	}
	if family == "http" {
		tabbar.Add(this.Lang(action, codes.Server_TabStat), "", "/servers/server/stat?serverId="+serverIdString, "chart area", selectedTabbar == "stat")
	}
	if family == "http" {
		tabbar.Add(this.Lang(action, codes.Server_TabAccessLogs), "", "/servers/server/log?serverId="+serverIdString, "history", selectedTabbar == "log")
	}
	tabbar.Add(this.Lang(action, codes.Server_TabSettings), "", "/servers/server/settings?serverId="+serverIdString, "setting", selectedTabbar == "setting")
	tabbar.Add(this.Lang(action, codes.Server_TabDelete), "", "/servers/server/delete?serverId="+serverIdString, "trash", selectedTabbar == "delete")

	actionutils.SetTabbar(action, tabbar)

	// 左侧操作子菜单
	switch types.String(mainTab) {
	case "board":
		action.Data["leftMenuItems"] = this.createBoardMenu(types.String(secondMenuItem), serverIdString, serverConfig, action)
	case "log":
		action.Data["leftMenuItems"] = this.createLogMenu(types.String(secondMenuItem), serverIdString, serverConfig, action)
	case "stat":
		action.Data["leftMenuItems"] = this.createStatMenu(types.String(secondMenuItem), serverIdString, serverConfig, action)
	case "setting":
		var menuItems = this.createSettingsMenu(types.String(secondMenuItem), serverIdString, serverConfig, action)
		action.Data["leftMenuItems"] = menuItems

		// 当前菜单
		action.Data["leftMenuActiveItem"] = nil
		for _, item := range menuItems {
			if item.GetBool("isActive") {
				action.Data["leftMenuActiveItem"] = item
				break
			}
		}
	case "delete":
		action.Data["leftMenuItems"] = this.createDeleteMenu(types.String(secondMenuItem), serverIdString, serverConfig, action)
	}
}

// 看板菜单
func (this *ServerHelper) createBoardMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig, actionPtr actions.ActionWrapper) []maps.Map {
	menuItems := []maps.Map{}
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuDashboard),
		"url":      "/servers/server/board?serverId=" + serverIdString,
		"isActive": secondMenuItem == "index",
	})
	return menuItems
}

// 日志菜单
func (this *ServerHelper) createLogMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig, actionPtr actions.ActionWrapper) []maps.Map {
	menuItems := []maps.Map{}
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuAccesslogRealtime),
		"url":      "/servers/server/log?serverId=" + serverIdString,
		"isActive": secondMenuItem == "index",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuAccesslogToday),
		"url":      "/servers/server/log/today?serverId=" + serverIdString,
		"isActive": secondMenuItem == "today",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuAccesslogHistory),
		"url":      "/servers/server/log/history?serverId=" + serverIdString,
		"isActive": secondMenuItem == "history",
	})
	return menuItems
}

// 统计菜单
func (this *ServerHelper) createStatMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig, actionPtr actions.ActionWrapper) []maps.Map {
	var menuItems = []maps.Map{}
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuStatTraffic),
		"url":      "/servers/server/stat?serverId=" + serverIdString,
		"isActive": secondMenuItem == "index",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuStatRegions),
		"url":      "/servers/server/stat/regions?serverId=" + serverIdString,
		"isActive": secondMenuItem == "region",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuStatProviders),
		"url":      "/servers/server/stat/providers?serverId=" + serverIdString,
		"isActive": secondMenuItem == "provider",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuStatClients),
		"url":      "/servers/server/stat/clients?serverId=" + serverIdString,
		"isActive": secondMenuItem == "client",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuStatWAF),
		"url":      "/servers/server/stat/waf?serverId=" + serverIdString,
		"isActive": secondMenuItem == "waf",
	})
	return menuItems
}

// 设置菜单
func (this *ServerHelper) createSettingsMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig, actionPtr actions.ActionWrapper) (items []maps.Map) {
	var menuItems = []maps.Map{
		{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingBasic),
			"url":      "/servers/server/settings?serverId=" + serverIdString,
			"isActive": secondMenuItem == "basic",
			"isOff":    !serverConfig.IsOn,
		},
	}

	// HTTP
	if serverConfig.IsHTTPFamily() {
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingDomains),
			"url":      "/servers/server/settings/serverNames?serverId=" + serverIdString,
			"isActive": secondMenuItem == "serverName",
			"isOn":     len(serverConfig.ServerNames) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingDNS),
			"url":      "/servers/server/settings/dns?serverId=" + serverIdString,
			"isActive": secondMenuItem == "dns",
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingHTTP),
			"url":      "/servers/server/settings/http?serverId=" + serverIdString,
			"isActive": secondMenuItem == "http",
			"isOn":     (serverConfig.HTTP != nil && serverConfig.HTTP.IsOn && len(serverConfig.HTTP.Listen) > 0) || (serverConfig.Web != nil && serverConfig.Web.RedirectToHttps != nil && serverConfig.Web.RedirectToHttps.IsOn),
			"isOff":    serverConfig.HTTP != nil && !serverConfig.HTTP.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingHTTPS),
			"url":      "/servers/server/settings/https?serverId=" + serverIdString,
			"isActive": secondMenuItem == "https",
			"isOn":     serverConfig.HTTPS != nil && serverConfig.HTTPS.IsOn && len(serverConfig.HTTPS.Listen) > 0,
			"isOff":    serverConfig.HTTPS != nil && !serverConfig.HTTPS.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingOrigins),
			"url":        "/servers/server/settings/reverseProxy?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "reverseProxy",
			"isOn":       serverConfig.ReverseProxyRef != nil && serverConfig.ReverseProxyRef.IsOn,
			"configCode": serverconfigs.ConfigCodeReverseProxy,
		})

		menuItems = this.filterMenuItems(serverConfig, menuItems, serverIdString, secondMenuItem, actionPtr)

		menuItems = append(menuItems, maps.Map{
			"name":     "-",
			"url":      "",
			"isActive": false,
		})
		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingRedirects),
			"url":        "/servers/server/settings/redirects?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "redirects",
			"isOn":       serverConfig.Web != nil && len(serverConfig.Web.HostRedirects) > 0,
			"configCode": serverconfigs.ConfigCodeHostRedirects,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingLocations),
			"url":      "/servers/server/settings/locations?serverId=" + serverIdString,
			"isActive": secondMenuItem == "locations",
			"isOn":     serverConfig.Web != nil && len(serverConfig.Web.Locations) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingRewriteRules),
			"url":      "/servers/server/settings/rewrite?serverId=" + serverIdString,
			"isActive": secondMenuItem == "rewrite",
			"isOn":     serverConfig.Web != nil && len(serverConfig.Web.RewriteRefs) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingWAF),
			"url":        "/servers/server/settings/waf?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "waf",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.FirewallRef != nil && serverConfig.Web.FirewallRef.IsOn,
			"configCode": serverconfigs.ConfigCodeWAF,
		})
		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingCache),
			"url":        "/servers/server/settings/cache?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "cache",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.Cache != nil && serverConfig.Web.Cache.IsOn,
			"configCode": serverconfigs.ConfigCodeCache,
		})

		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingAuth),
			"url":        "/servers/server/settings/access?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "access",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.Auth != nil && serverConfig.Web.Auth.IsOn,
			"configCode": serverconfigs.ConfigCodeAuth,
		})

		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingReferers),
			"url":        "/servers/server/settings/referers?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "referer",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.Referers != nil && serverConfig.Web.Referers.IsOn,
			"configCode": serverconfigs.ConfigCodeReferers,
		})
		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingUserAgents),
			"url":        "/servers/server/settings/userAgent?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "userAgent",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.UserAgent != nil && serverConfig.Web.UserAgent.IsOn,
			"configCode": serverconfigs.ConfigCodeUserAgent,
		})

		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingAccessLog),
			"url":        "/servers/server/settings/accessLog?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "accessLog",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.AccessLogRef != nil && serverConfig.Web.AccessLogRef.IsOn,
			"configCode": serverconfigs.ConfigCodeAccessLog,
		})

		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingCompression),
			"url":        "/servers/server/settings/compression?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "compression",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.Compression != nil && serverConfig.Web.Compression.IsOn,
			"configCode": serverconfigs.ConfigCodeCompression,
		})

		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingPages),
			"url":        "/servers/server/settings/pages?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "pages",
			"isOn":       serverConfig.Web != nil && (len(serverConfig.Web.Pages) > 0 || (serverConfig.Web.Shutdown != nil && serverConfig.Web.Shutdown.IsOn)),
			"configCode": serverconfigs.ConfigCodePages,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingHTTPHeaders),
			"url":      "/servers/server/settings/headers?serverId=" + serverIdString,
			"isActive": secondMenuItem == "header",
			"isOn":     this.hasHTTPHeaders(serverConfig.Web),
		})
		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingWebsocket),
			"url":        "/servers/server/settings/websocket?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "websocket",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.WebsocketRef != nil && serverConfig.Web.WebsocketRef.IsOn,
			"configCode": serverconfigs.ConfigCodeWebsocket,
		})
		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingWebP),
			"url":        "/servers/server/settings/webp?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "webp",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.WebP != nil && serverConfig.Web.WebP.IsOn,
			"configCode": serverconfigs.ConfigCodeWebp,
		})

		menuItems = this.filterMenuItems3(serverConfig, menuItems, serverIdString, secondMenuItem, actionPtr)

		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingStat),
			"url":        "/servers/server/settings/stat?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "stat",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.StatRef != nil && serverConfig.Web.StatRef.IsOn,
			"configCode": serverconfigs.ConfigCodeStat,
		})
		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingCharset),
			"url":        "/servers/server/settings/charset?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "charset",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.Charset != nil && serverConfig.Web.Charset.IsOn,
			"configCode": serverconfigs.ConfigCodeCharset,
		})
		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingRoot),
			"url":        "/servers/server/settings/web?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "web",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.Root != nil && serverConfig.Web.Root.IsOn,
			"configCode": serverconfigs.ConfigCodeRoot,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingFastcgi),
			"url":      "/servers/server/settings/fastcgi?serverId=" + serverIdString,
			"isActive": secondMenuItem == "fastcgi",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.FastcgiRef != nil && serverConfig.Web.FastcgiRef.IsOn,
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "-",
			"url":      "",
			"isActive": false,
		})

		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingClientIP),
			"url":        "/servers/server/settings/remoteAddr?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "remoteAddr",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.RemoteAddr != nil && serverConfig.Web.RemoteAddr.IsOn,
			"configCode": serverconfigs.ConfigCodeRemoteAddr,
		})

		menuItems = append(menuItems, maps.Map{
			"name":       this.Lang(actionPtr, codes.Server_MenuSettingRequestLimit),
			"url":        "/servers/server/settings/requestLimit?serverId=" + serverIdString,
			"isActive":   secondMenuItem == "requestLimit",
			"isOn":       serverConfig.Web != nil && serverConfig.Web.RequestLimit != nil && serverConfig.Web.RequestLimit.IsOn,
			"configCode": serverconfigs.ConfigCodeRequestLimit,
		})

		menuItems = this.filterMenuItems2(serverConfig, menuItems, serverIdString, secondMenuItem, actionPtr)

		menuItems = append(menuItems, maps.Map{
			"name":     "-",
			"url":      "",
			"isActive": false,
		})

		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingOthers),
			"url":      "/servers/server/settings/common?serverId=" + serverIdString,
			"isActive": secondMenuItem == "common",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.MergeSlashes,
		})
	} else if serverConfig.IsTCPFamily() {
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingDNS),
			"url":      "/servers/server/settings/dns?serverId=" + serverIdString,
			"isActive": secondMenuItem == "dns",
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingTCP),
			"url":      "/servers/server/settings/tcp?serverId=" + serverIdString,
			"isActive": secondMenuItem == "tcp",
			"isOn":     serverConfig.TCP != nil && serverConfig.TCP.IsOn && len(serverConfig.TCP.Listen) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingTLS),
			"url":      "/servers/server/settings/tls?serverId=" + serverIdString,
			"isActive": secondMenuItem == "tls",
			"isOn":     serverConfig.TLS != nil && serverConfig.TLS.IsOn && len(serverConfig.TLS.Listen) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingOrigins),
			"url":      "/servers/server/settings/reverseProxy?serverId=" + serverIdString,
			"isActive": secondMenuItem == "reverseProxy",
			"isOn":     serverConfig.ReverseProxyRef != nil && serverConfig.ReverseProxyRef.IsOn,
		})
	} else if serverConfig.IsUDPFamily() {
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingDNS),
			"url":      "/servers/server/settings/dns?serverId=" + serverIdString,
			"isActive": secondMenuItem == "dns",
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingUDP),
			"url":      "/servers/server/settings/udp?serverId=" + serverIdString,
			"isActive": secondMenuItem == "udp",
			"isOn":     serverConfig.UDP != nil && serverConfig.UDP.IsOn && len(serverConfig.UDP.Listen) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     this.Lang(actionPtr, codes.Server_MenuSettingOrigins),
			"url":      "/servers/server/settings/reverseProxy?serverId=" + serverIdString,
			"isActive": secondMenuItem == "reverseProxy",
			"isOn":     serverConfig.ReverseProxyRef != nil && serverConfig.ReverseProxyRef.IsOn,
		})
	}

	return menuItems
}

// 删除菜单
func (this *ServerHelper) createDeleteMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig, actionPtr actions.ActionWrapper) []maps.Map {
	menuItems := []maps.Map{}
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingDelete),
		"url":      "/servers/server/delete?serverId=" + serverIdString,
		"isActive": secondMenuItem == "index",
	})
	return menuItems
}

// 检查是否已设置Header
func (this *ServerHelper) hasHTTPHeaders(web *serverconfigs.HTTPWebConfig) bool {
	if web == nil {
		return false
	}
	if web.RequestHeaderPolicyRef != nil {
		if web.RequestHeaderPolicyRef.IsOn && web.RequestHeaderPolicy != nil && !web.RequestHeaderPolicy.IsEmpty() {
			return true
		}
	}
	if web.ResponseHeaderPolicyRef != nil {
		if web.ResponseHeaderPolicyRef.IsOn && web.ResponseHeaderPolicy != nil && !web.ResponseHeaderPolicy.IsEmpty() {
			return true
		}
	}
	return false
}
