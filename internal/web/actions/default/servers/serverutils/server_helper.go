package serverutils

import (
	"encoding/json"
	"errors"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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
	mainTab, _ := action.Data["mainTab"]
	secondMenuItem, _ := action.Data["secondMenuItem"]

	serverId := action.ParamInt64("serverId")
	if serverId == 0 {
		return
	}
	serverIdString := strconv.FormatInt(serverId, 10)
	action.Data["serverId"] = serverId

	// 读取server信息
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		logs.Error(err)
		return
	}

	serverResp, err := rpcClient.ServerRPC().FindEnabledServer(rpcClient.Context(action.Context.GetInt64("adminId")), &pb.FindEnabledServerRequest{ServerId: serverId})
	if err != nil {
		logs.Error(err)
		return
	}
	server := serverResp.Server
	if server == nil {
		logs.Error(errors.New("can not find the server"))
		return
	}

	// 服务管理
	serverConfig := &serverconfigs.ServerConfig{}
	err = json.Unmarshal(server.Config, serverConfig)
	if err != nil {
		logs.Error(err)
		return
	}

	// 协议簇
	family := ""
	if serverConfig.IsHTTPFamily() {
		family = "http"
	} else if serverConfig.IsTCPFamily() {
		family = "tcp"
	} else if serverConfig.IsUnixFamily() {
		family = "unix"
	} else if serverConfig.IsUDPFamily() {
		family = "udp"
	}
	action.Data["serverFamily"] = family

	// TABBAR
	selectedTabbar, _ := action.Data["mainTab"]
	tabbar := actionutils.NewTabbar()
	tabbar.Add("服务列表", "", "/servers", "", false)
	if teaconst.IsPlus {
		tabbar.Add("看板", "", "/servers/server/boards?serverId="+serverIdString, "dashboard", selectedTabbar == "board")
	}
	if family == "http" {
		tabbar.Add("统计", "", "/servers/server/stat?serverId="+serverIdString, "chart area", selectedTabbar == "stat")
	}
	if family == "http" {
		tabbar.Add("日志", "", "/servers/server/log?serverId="+serverIdString, "history", selectedTabbar == "log")
	}
	tabbar.Add("设置", "", "/servers/server/settings?serverId="+serverIdString, "setting", selectedTabbar == "setting")
	tabbar.Add("删除", "", "/servers/server/delete?serverId="+serverIdString, "trash", selectedTabbar == "delete")
	{
		m := tabbar.Add("当前服务："+server.Name, "", "/servers/server?serverId="+serverIdString, "", false)
		m["right"] = true
	}

	actionutils.SetTabbar(action, tabbar)

	// 左侧操作子菜单
	switch types.String(mainTab) {
	case "board":
		action.Data["leftMenuItems"] = this.createBoardMenu(types.String(secondMenuItem), serverIdString, serverConfig)
	case "log":
		action.Data["leftMenuItems"] = this.createLogMenu(types.String(secondMenuItem), serverIdString, serverConfig)
	case "stat":
		action.Data["leftMenuItems"] = this.createStatMenu(types.String(secondMenuItem), serverIdString, serverConfig)
	case "setting":
		action.Data["leftMenuItems"] = this.createSettingsMenu(types.String(secondMenuItem), serverIdString, serverConfig)
	case "delete":
		action.Data["leftMenuItems"] = this.createDeleteMenu(types.String(secondMenuItem), serverIdString, serverConfig)
	}
}

// 看板菜单
func (this *ServerHelper) createBoardMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig) []maps.Map {
	menuItems := []maps.Map{}
	menuItems = append(menuItems, maps.Map{
		"name":     "看板",
		"url":      "/servers/server/board?serverId=" + serverIdString,
		"isActive": secondMenuItem == "index",
	})
	return menuItems
}

// 日志菜单
func (this *ServerHelper) createLogMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig) []maps.Map {
	menuItems := []maps.Map{}
	menuItems = append(menuItems, maps.Map{
		"name":     "实时",
		"url":      "/servers/server/log?serverId=" + serverIdString,
		"isActive": secondMenuItem == "index",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     "今天",
		"url":      "/servers/server/log/today?serverId=" + serverIdString,
		"isActive": secondMenuItem == "today",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     "历史",
		"url":      "/servers/server/log/history?serverId=" + serverIdString,
		"isActive": secondMenuItem == "history",
	})
	return menuItems
}

// 统计菜单
func (this *ServerHelper) createStatMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig) []maps.Map {
	menuItems := []maps.Map{}
	menuItems = append(menuItems, maps.Map{
		"name":     "流量统计",
		"url":      "/servers/server/stat?serverId=" + serverIdString,
		"isActive": secondMenuItem == "index",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     "地域分布",
		"url":      "/servers/server/stat/regions?serverId=" + serverIdString,
		"isActive": secondMenuItem == "region",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     "运营商",
		"url":      "/servers/server/stat/providers?serverId=" + serverIdString,
		"isActive": secondMenuItem == "provider",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     "终端",
		"url":      "/servers/server/stat/clients?serverId=" + serverIdString,
		"isActive": secondMenuItem == "client",
	})
	menuItems = append(menuItems, maps.Map{
		"name":     "WAF",
		"url":      "/servers/server/stat/waf?serverId=" + serverIdString,
		"isActive": secondMenuItem == "waf",
	})
	return menuItems
}

// 设置菜单
func (this *ServerHelper) createSettingsMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig) (items []maps.Map) {
	menuItems := []maps.Map{
		{
			"name":     "基本信息",
			"url":      "/servers/server/settings?serverId=" + serverIdString,
			"isActive": secondMenuItem == "basic",
			"isOff":    !serverConfig.IsOn,
		},
		{
			"name":     "DNS",
			"url":      "/servers/server/settings/dns?serverId=" + serverIdString,
			"isActive": secondMenuItem == "dns",
		},
	}

	// HTTP
	if serverConfig.IsHTTPFamily() {
		menuItems = append(menuItems, maps.Map{
			"name":     "域名",
			"url":      "/servers/server/settings/serverNames?serverId=" + serverIdString,
			"isActive": secondMenuItem == "serverName",
			"isOn":     len(serverConfig.ServerNames) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "HTTP",
			"url":      "/servers/server/settings/http?serverId=" + serverIdString,
			"isActive": secondMenuItem == "http",
			"isOn":     (serverConfig.HTTP != nil && serverConfig.HTTP.IsOn && len(serverConfig.HTTP.Listen) > 0) || (serverConfig.Web != nil && serverConfig.Web.RedirectToHttps != nil && serverConfig.Web.RedirectToHttps.IsOn),
			"isOff":    serverConfig.HTTP != nil && !serverConfig.HTTP.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "HTTPS",
			"url":      "/servers/server/settings/https?serverId=" + serverIdString,
			"isActive": secondMenuItem == "https",
			"isOn":     serverConfig.HTTPS != nil && serverConfig.HTTPS.IsOn && len(serverConfig.HTTPS.Listen) > 0,
			"isOff":    serverConfig.HTTPS != nil && !serverConfig.HTTPS.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "Web设置",
			"url":      "/servers/server/settings/web?serverId=" + serverIdString,
			"isActive": secondMenuItem == "web",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.Root != nil && serverConfig.Web.Root.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "反向代理",
			"url":      "/servers/server/settings/reverseProxy?serverId=" + serverIdString,
			"isActive": secondMenuItem == "reverseProxy",
			"isOn":     serverConfig.ReverseProxyRef != nil && serverConfig.ReverseProxyRef.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "-",
			"url":      "",
			"isActive": false,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "URL跳转",
			"url":      "/servers/server/settings/redirects?serverId=" + serverIdString,
			"isActive": secondMenuItem == "redirects",
			"isOn":     serverConfig.Web != nil && len(serverConfig.Web.HostRedirects) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "路由规则",
			"url":      "/servers/server/settings/locations?serverId=" + serverIdString,
			"isActive": secondMenuItem == "locations",
			"isOn":     serverConfig.Web != nil && len(serverConfig.Web.Locations) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "重写规则",
			"url":      "/servers/server/settings/rewrite?serverId=" + serverIdString,
			"isActive": secondMenuItem == "rewrite",
			"isOn":     serverConfig.Web != nil && len(serverConfig.Web.RewriteRefs) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "WAF",
			"url":      "/servers/server/settings/waf?serverId=" + serverIdString,
			"isActive": secondMenuItem == "waf",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.FirewallRef != nil && serverConfig.Web.FirewallRef.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "缓存",
			"url":      "/servers/server/settings/cache?serverId=" + serverIdString,
			"isActive": secondMenuItem == "cache",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.Cache != nil && serverConfig.Web.Cache.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "访问控制",
			"url":      "/servers/server/settings/access?serverId=" + serverIdString,
			"isActive": secondMenuItem == "access",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.Auth != nil && serverConfig.Web.Auth.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "字符编码",
			"url":      "/servers/server/settings/charset?serverId=" + serverIdString,
			"isActive": secondMenuItem == "charset",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.Charset != nil && serverConfig.Web.Charset.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "访问日志",
			"url":      "/servers/server/settings/accessLog?serverId=" + serverIdString,
			"isActive": secondMenuItem == "accessLog",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.AccessLogRef != nil && serverConfig.Web.AccessLogRef.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "统计",
			"url":      "/servers/server/settings/stat?serverId=" + serverIdString,
			"isActive": secondMenuItem == "stat",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.StatRef != nil && serverConfig.Web.StatRef.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "内容压缩",
			"url":      "/servers/server/settings/compression?serverId=" + serverIdString,
			"isActive": secondMenuItem == "compression",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.Compression != nil && serverConfig.Web.Compression.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "特殊页面",
			"url":      "/servers/server/settings/pages?serverId=" + serverIdString,
			"isActive": secondMenuItem == "pages",
			"isOn":     serverConfig.Web != nil && (len(serverConfig.Web.Pages) > 0 || (serverConfig.Web.Shutdown != nil && serverConfig.Web.Shutdown.IsOn)),
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "HTTP Header",
			"url":      "/servers/server/settings/headers?serverId=" + serverIdString,
			"isActive": secondMenuItem == "header",
			"isOn":     this.hasHTTPHeaders(serverConfig.Web),
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "Websocket",
			"url":      "/servers/server/settings/websocket?serverId=" + serverIdString,
			"isActive": secondMenuItem == "websocket",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.WebsocketRef != nil && serverConfig.Web.WebsocketRef.IsOn,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "WebP",
			"url":      "/servers/server/settings/webp?serverId=" + serverIdString,
			"isActive": secondMenuItem == "webp",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.WebP != nil && serverConfig.Web.WebP.IsOn,
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "Fastcgi",
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
			"name":     "访客IP地址",
			"url":      "/servers/server/settings/remoteAddr?serverId=" + serverIdString,
			"isActive": secondMenuItem == "remoteAddr",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.RemoteAddr != nil && serverConfig.Web.RemoteAddr.IsOn,
		})

		if teaconst.IsPlus {
			menuItems = append(menuItems, maps.Map{
				"name":     "流量限制",
				"url":      "/servers/server/settings/traffic?serverId=" + serverIdString,
				"isActive": secondMenuItem == "traffic",
				"isOn":     serverConfig.TrafficLimit != nil && serverConfig.TrafficLimit.IsOn,
			})
		}

		menuItems = append(menuItems, maps.Map{
			"name":     "-",
			"url":      "",
			"isActive": false,
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "其他设置",
			"url":      "/servers/server/settings/common?serverId=" + serverIdString,
			"isActive": secondMenuItem == "common",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.MergeSlashes,
		})
	} else if serverConfig.IsTCPFamily() {
		menuItems = append(menuItems, maps.Map{
			"name":     "TCP",
			"url":      "/servers/server/settings/tcp?serverId=" + serverIdString,
			"isActive": secondMenuItem == "tcp",
			"isOn":     serverConfig.TCP != nil && serverConfig.TCP.IsOn && len(serverConfig.TCP.Listen) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "TLS",
			"url":      "/servers/server/settings/tls?serverId=" + serverIdString,
			"isActive": secondMenuItem == "tls",
			"isOn":     serverConfig.TLS != nil && serverConfig.TLS.IsOn && len(serverConfig.TLS.Listen) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "反向代理",
			"url":      "/servers/server/settings/reverseProxy?serverId=" + serverIdString,
			"isActive": secondMenuItem == "reverseProxy",
			"isOn":     serverConfig.ReverseProxyRef != nil && serverConfig.ReverseProxyRef.IsOn,
		})
	} else if serverConfig.IsUnixFamily() {
		menuItems = append(menuItems, maps.Map{
			"name":     "Unix",
			"url":      "/servers/server/settings/unix?serverId=" + serverIdString,
			"isActive": secondMenuItem == "unix",
			"isOn":     serverConfig.Unix != nil && serverConfig.Unix.IsOn && len(serverConfig.Unix.Listen) > 0,
		})
	} else if serverConfig.IsUDPFamily() {
		menuItems = append(menuItems, maps.Map{
			"name":     "UDP",
			"url":      "/servers/server/settings/udp?serverId=" + serverIdString,
			"isActive": secondMenuItem == "udp",
			"isOn":     serverConfig.UDP != nil && serverConfig.UDP.IsOn && len(serverConfig.UDP.Listen) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "反向代理",
			"url":      "/servers/server/settings/reverseProxy?serverId=" + serverIdString,
			"isActive": secondMenuItem == "reverseProxy",
			"isOn":     serverConfig.ReverseProxyRef != nil && serverConfig.ReverseProxyRef.IsOn,
		})
	}

	return menuItems
}

// 删除菜单
func (this *ServerHelper) createDeleteMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig) []maps.Map {
	menuItems := []maps.Map{}
	menuItems = append(menuItems, maps.Map{
		"name":     "删除",
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
