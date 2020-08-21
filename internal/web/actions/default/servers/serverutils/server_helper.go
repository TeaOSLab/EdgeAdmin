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
	tabbar.Add("当前："+serverConfig.Name, "", "/servers", "left long alternate arrow", false)
	tabbar.Add("看板", "", "/servers/server/board?serverId="+serverIdString, "dashboard", selectedTabbar == "board")
	tabbar.Add("日志", "", "/servers/server/log?serverId="+serverIdString, "history", selectedTabbar == "log")
	tabbar.Add("统计", "", "/servers/server/stat?serverId="+serverIdString, "chart area", selectedTabbar == "stat")
	tabbar.Add("设置", "", "/servers/server/settings?serverId="+serverIdString, "setting", selectedTabbar == "setting")
	tabbar.Add("删除", "", "/servers/server/delete?serverId="+serverIdString, "trash", selectedTabbar == "delete")

	actionutils.SetTabbar(action, tabbar)

	// 左侧操作子菜单
	switch types.String(mainTab) {
	case "board":
		// TODO
	case "log":
		// TODO
	case "stat":
		// TODO
	case "setting":
		action.Data["leftMenuItems"] = this.createSettingsMenu(types.String(secondMenuItem), serverIdString, serverConfig)
	case "delete":
		// TODO
	}
}

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
			"url":      "/servers/server/http?serverId=" + serverIdString,
			"isActive": secondMenuItem == "http",
		})
		menuItems = append(menuItems, maps.Map{
			"name":     "HTTPS",
			"url":      "/servers/server/https?serverId=" + serverIdString,
			"isActive": secondMenuItem == "https",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "Web设置",
			"url":      "/servers/server/web?serverId=" + serverIdString,
			"isActive": secondMenuItem == "web",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "字符集",
			"url":      "/servers/server/charset?serverId=" + serverIdString,
			"isActive": secondMenuItem == "charset",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "访问日志",
			"url":      "/servers/server/accessLog?serverId=" + serverIdString,
			"isActive": secondMenuItem == "accessLog",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "统计",
			"url":      "/servers/server/stat?serverId=" + serverIdString,
			"isActive": secondMenuItem == "stat",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "Gzip压缩",
			"url":      "/servers/server/gzip?serverId=" + serverIdString,
			"isActive": secondMenuItem == "gzip",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "特殊页面",
			"url":      "/servers/server/pages?serverId=" + serverIdString,
			"isActive": secondMenuItem == "pages",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "HTTP Header",
			"url":      "/servers/server/headers?serverId=" + serverIdString,
			"isActive": secondMenuItem == "header",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "反向代理",
			"url":      "/servers/server/reverseProxy?serverId=" + serverIdString,
			"isActive": secondMenuItem == "reverseProxy",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "路径规则",
			"url":      "/servers/server/locations?serverId=" + serverIdString,
			"isActive": secondMenuItem == "locations",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "访问控制",
			"url":      "/servers/server/access?serverId=" + serverIdString,
			"isActive": secondMenuItem == "access",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "WAF",
			"url":      "/servers/server/waf?serverId=" + serverIdString,
			"isActive": secondMenuItem == "waf",
		})

		menuItems = append(menuItems, maps.Map{
			"name":     "缓存",
			"url":      "/servers/server/cache?serverId=" + serverIdString,
			"isActive": secondMenuItem == "cache",
		})
	} else if serverConfig.IsTCP() {
		menuItems = append(menuItems, maps.Map{
			"name":     "TCP",
			"url":      "/servers/server/tcp?serverId=" + serverIdString,
			"isActive": secondMenuItem == "tcp",
		})
	} else if serverConfig.IsUnix() {
		menuItems = append(menuItems, maps.Map{
			"name":     "Unix",
			"url":      "/servers/server/unix?serverId=" + serverIdString,
			"isActive": secondMenuItem == "unix",
		})
	} else if serverConfig.IsUDP() {
		menuItems = append(menuItems, maps.Map{
			"name":     "UDP",
			"url":      "/servers/server/udp?serverId=" + serverIdString,
			"isActive": secondMenuItem == "udp",
		})
	}

	return menuItems
}
