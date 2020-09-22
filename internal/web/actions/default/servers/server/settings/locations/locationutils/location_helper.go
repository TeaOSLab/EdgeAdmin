package locationutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"reflect"
)

type LocationHelper struct {
}

func NewLocationHelper() *LocationHelper {
	return &LocationHelper{}
}

func (this *LocationHelper) BeforeAction(actionPtr actions.ActionWrapper) {
	action := actionPtr.Object()
	if action.Request.Method != http.MethodGet {
		return
	}

	serverIdString := action.ParamString("serverId")
	locationIdString := action.ParamString("locationId")

	action.Data["leftMenuItemIsDisabled"] = true
	action.Data["mainMenu"] = "server"
	action.Data["mainTab"] = "setting"
	action.Data["secondMenuItem"] = "locations"
	action.Data["tinyLeftMenuItems"] = this.createMenus(serverIdString, locationIdString, action.Data.GetString("tinyMenuItem"))

	// 路径信息
	parentActionValue := reflect.ValueOf(actionPtr).Elem().FieldByName("ParentAction")
	if parentActionValue.IsValid() {
		parentAction, isOk := parentActionValue.Interface().(actionutils.ParentAction)
		if isOk {
			locationId := action.ParamInt64("locationId")
			locationConfig, isOk := FindLocationConfig(&parentAction, locationId)
			if !isOk {
				return
			}
			action.Data["locationId"] = locationId
			action.Data["locationConfig"] = locationConfig
		}
	}
}

func (this *LocationHelper) createMenus(serverIdString string, locationIdString string, secondMenuItem string) []maps.Map {
	menuItems := []maps.Map{
		{
			"name":     "基本信息",
			"url":      "/servers/server/settings/locations/location?serverId=" + serverIdString + "&locationId=" + locationIdString,
			"isActive": secondMenuItem == "basic",
		},
	}

	menuItems = append(menuItems, maps.Map{
		"name":     "Web设置",
		"url":      "/servers/server/settings/locations/web?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "web",
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "反向代理",
		"url":      "/servers/server/settings/locations/reverseProxy?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "reverseProxy",
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "访问控制",
		"url":      "/servers/server/settings/locations/access?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "access",
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "WAF",
		"url":      "/servers/server/settings/locations/waf?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "waf",
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "缓存",
		"url":      "/servers/server/settings/locations/cache?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "cache",
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "-",
		"url":      "",
		"isActive": false,
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "字符编码",
		"url":      "/servers/server/settings/locations/charset?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "charset",
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "访问日志",
		"url":      "/servers/server/settings/locations/accessLog?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "accessLog",
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "统计",
		"url":      "/servers/server/settings/locations/stat?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "stat",
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "Gzip压缩",
		"url":      "/servers/server/settings/locations/gzip?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "gzip",
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "特殊页面",
		"url":      "/servers/server/settings/locations/pages?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "pages",
	})

	menuItems = append(menuItems, maps.Map{
		"name":     "HTTP Header",
		"url":      "/servers/server/settings/locations/headers?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "header",
	})

	return menuItems
}
