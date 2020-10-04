package serverutils

import (
	"encoding/json"
	"errors"
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

	// TABBAR
	selectedTabbar, _ := action.Data["mainTab"]
	tabbar := actionutils.NewTabbar()
	tabbar.Add("服务列表", "", "/servers", "", false)
	//tabbar.Add("看板", "", "/servers/server/board?serverId="+serverIdString, "dashboard", selectedTabbar == "board")
	tabbar.Add("日志", "", "/servers/server/log?serverId="+serverIdString, "history", selectedTabbar == "log")
	//tabbar.Add("统计", "", "/servers/server/stat?serverId="+serverIdString, "chart area", selectedTabbar == "stat")
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
	return menuItems
}

// 统计菜单
func (this *ServerHelper) createStatMenu(secondMenuItem string, serverIdString string, serverConfig *serverconfigs.ServerConfig) []maps.Map {
	menuItems := []maps.Map{}
	menuItems = append(menuItems, maps.Map{
		"name":     "统计",
		"url":      "/servers/server/stat?serverId=" + serverIdString,
		"isActive": secondMenuItem == "index",
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
	}

	// HTTP
	if serverConfig.IsHTTP() {
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
			"name":     "路径规则",
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
			"isOn":     serverConfig.Web != nil && len(serverConfig.Web.CacheRefs) > 0,
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "访问控制",
			"url":      "/servers/server/settings/access?serverId=" + serverIdString,
			"isActive": secondMenuItem == "access",
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
			"name":     "Gzip压缩",
			"url":      "/servers/server/settings/gzip?serverId=" + serverIdString,
			"isActive": secondMenuItem == "gzip",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.GzipRef != nil && serverConfig.Web.GzipRef.IsOn,
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
			"isOn":     serverConfig.Web != nil && ((serverConfig.Web.RequestHeaderPolicyRef != nil && serverConfig.Web.RequestHeaderPolicyRef.IsOn) || (serverConfig.Web.ResponseHeaderPolicyRef != nil && serverConfig.Web.ResponseHeaderPolicyRef.IsOn)),
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "Websocket",
			"url":      "/servers/server/settings/websocket?serverId=" + serverIdString,
			"isActive": secondMenuItem == "websocket",
			"isOn":     serverConfig.Web != nil && serverConfig.Web.WebsocketRef != nil && serverConfig.Web.WebsocketRef.IsOn,
		})
	} else if serverConfig.IsTCP() {
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
	} else if serverConfig.IsUnix() {
		menuItems = append(menuItems, maps.Map{
			"name":     "Unix",
			"url":      "/servers/server/settings/unix?serverId=" + serverIdString,
			"isActive": secondMenuItem == "unix",
			"isOn":     serverConfig.Unix != nil && serverConfig.Unix.IsOn && len(serverConfig.Unix.Listen) > 0,
		})
	} else if serverConfig.IsUDP() {
		menuItems = append(menuItems, maps.Map{
			"name":     "UDP",
			"url":      "/servers/server/settings/udp?serverId=" + serverIdString,
			"isActive": secondMenuItem == "udp",
			"isOn":     serverConfig.UDP != nil && serverConfig.UDP.IsOn && len(serverConfig.UDP.Listen) > 0,
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
