package serverutils

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs/serverconfigs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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
	action.Data["leftMenuItems"] = []maps.Map{}
	mainTab, _ := action.Data["mainTab"]
	secondMenuItem, _ := action.Data["secondMenuItem"]

	serverId := action.ParamInt64("serverId")
	serverIdString := strconv.FormatInt(serverId, 10)

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

	// 源站管理
	serverConfig := &serverconfigs.ServerConfig{}
	err = json.Unmarshal(server.Config, serverConfig)
	if err != nil {
		logs.Error(err)
		return
	}

	// TABBAR
	selectedTabbar, _ := action.Data["mainTab"]
	tabbar := actionutils.NewTabbar()
	tabbar.Add("当前服务："+serverConfig.Name, "", "/servers", "left long alternate arrow", false)
	tabbar.Add("看板", "", "/servers/server/board?serverId="+serverIdString, "dashboard", selectedTabbar == "board")
	tabbar.Add("日志", "", "/servers/server/log?serverId="+serverIdString, "history", selectedTabbar == "log")
	tabbar.Add("统计", "", "/servers/server/stat?serverId="+serverIdString, "chart area", selectedTabbar == "stat")
	tabbar.Add("设置", "", "/servers/server/settings?serverId="+serverIdString, "setting", selectedTabbar == "setting")
	tabbar.Add("删除", "", "/servers/server/delete?serverId="+serverIdString, "trash", selectedTabbar == "delete")

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
		},
	}

	// HTTP
	if serverConfig.IsHTTP() {
		menuItems = append(menuItems, maps.Map{
			"name":     "HTTP",
			"url":      "/servers/server/settings/http?serverId=" + serverIdString,
			"isActive": secondMenuItem == "http",
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "HTTPS",
			"url":      "/servers/server/settings/https?serverId=" + serverIdString,
			"isActive": secondMenuItem == "https",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "Web设置",
			"url":      "/servers/server/settings/web?serverId=" + serverIdString,
			"isActive": secondMenuItem == "web",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "反向代理",
			"url":      "/servers/server/settings/reverseProxy?serverId=" + serverIdString,
			"isActive": secondMenuItem == "reverseProxy",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "路径规则",
			"url":      "/servers/server/settings/locations?serverId=" + serverIdString,
			"isActive": secondMenuItem == "locations",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "访问控制",
			"url":      "/servers/server/settings/access?serverId=" + serverIdString,
			"isActive": secondMenuItem == "access",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "WAF",
			"url":      "/servers/server/settings/waf?serverId=" + serverIdString,
			"isActive": secondMenuItem == "waf",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "缓存",
			"url":      "/servers/server/settings/cache?serverId=" + serverIdString,
			"isActive": secondMenuItem == "cache",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "-",
			"url":      "",
			"isActive": false,
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "字符集",
			"url":      "/servers/server/settings/charset?serverId=" + serverIdString,
			"isActive": secondMenuItem == "charset",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "访问日志",
			"url":      "/servers/server/settings/accessLog?serverId=" + serverIdString,
			"isActive": secondMenuItem == "accessLog",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "统计",
			"url":      "/servers/server/settings/stat?serverId=" + serverIdString,
			"isActive": secondMenuItem == "stat",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "Gzip压缩",
			"url":      "/servers/server/settings/gzip?serverId=" + serverIdString,
			"isActive": secondMenuItem == "gzip",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "特殊页面",
			"url":      "/servers/server/settings/pages?serverId=" + serverIdString,
			"isActive": secondMenuItem == "pages",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "HTTP Header",
			"url":      "/servers/server/settings/headers?serverId=" + serverIdString,
			"isActive": secondMenuItem == "header",
		})
	} else if serverConfig.IsTCP() {
		menuItems = append(menuItems, maps.Map{
			"name":     "TCP",
			"url":      "/servers/server/settings/tcp?serverId=" + serverIdString,
			"isActive": secondMenuItem == "tcp",
		})
	} else if serverConfig.IsUnix() {
		menuItems = append(menuItems, maps.Map{
			"name":     "Unix",
			"url":      "/servers/server/settings/unix?serverId=" + serverIdString,
			"isActive": secondMenuItem == "unix",
		})
	} else if serverConfig.IsUDP() {
		menuItems = append(menuItems, maps.Map{
			"name":     "UDP",
			"url":      "/servers/server/settings/udp?serverId=" + serverIdString,
			"isActive": secondMenuItem == "udp",
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
