package locationutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/http"
)

type LocationHelper struct {
	helpers.LangHelper
}

func NewLocationHelper() *LocationHelper {
	return &LocationHelper{}
}

func (this *LocationHelper) BeforeAction(actionPtr actions.ActionWrapper) {
	var action = actionPtr.Object()
	if action.Request.Method != http.MethodGet {
		return
	}

	var serverIdString = action.ParamString("serverId")
	var locationIdString = action.ParamString("locationId")

	action.Data["leftMenuItemIsDisabled"] = true
	action.Data["mainMenu"] = "server"
	action.Data["mainTab"] = "setting"
	action.Data["secondMenuItem"] = "locations"

	// 路径信息
	var currentLocationConfig *serverconfigs.HTTPLocationConfig = nil
	parentActionValue, ok := actionPtr.(interface {
		Parent() *actionutils.ParentAction
	})
	if ok {
		var parentAction = parentActionValue.Parent()
		var locationId = action.ParamInt64("locationId")
		locationConfig, isOk := FindLocationConfig(parentAction, locationId)
		if !isOk {
			return
		}
		action.Data["locationId"] = locationId
		action.Data["locationConfig"] = locationConfig
		currentLocationConfig = locationConfig
	}

	// 左侧菜单
	action.Data["tinyLeftMenuItems"] = this.createMenus(serverIdString, locationIdString, action.Data.GetString("tinyMenuItem"), currentLocationConfig, actionPtr)
}

func (this *LocationHelper) createMenus(serverIdString string, locationIdString string, secondMenuItem string, locationConfig *serverconfigs.HTTPLocationConfig, actionPtr actions.ActionWrapper) []maps.Map {
	menuItems := []maps.Map{}
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingBasic),
		"url":      "/servers/server/settings/locations/location?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "basic",
		"isOff":    locationConfig != nil && !locationConfig.IsOn,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingHTTP),
		"url":      "/servers/server/settings/locations/http?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "http",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.RedirectToHttps != nil && locationConfig.Web.RedirectToHttps.IsPrior,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingOrigins),
		"url":      "/servers/server/settings/locations/reverseProxy?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "reverseProxy",
		"isOn":     locationConfig != nil && locationConfig.ReverseProxyRef != nil && locationConfig.ReverseProxyRef.IsPrior,
	})

	menuItems = this.filterMenuItems1(locationConfig, menuItems, serverIdString, locationIdString, secondMenuItem, actionPtr)

	menuItems = append(menuItems, maps.Map{
		"name":     "-",
		"url":      "",
		"isActive": false,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingRewriteRules),
		"url":      "/servers/server/settings/locations/rewrite?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "rewrite",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && len(locationConfig.Web.RewriteRefs) > 0,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingWAF),
		"url":      "/servers/server/settings/locations/waf?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "waf",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.FirewallRef != nil && locationConfig.Web.FirewallRef.IsPrior,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingCache),
		"url":      "/servers/server/settings/locations/cache?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "cache",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.Cache != nil && locationConfig.Web.Cache.IsPrior && locationConfig.Web.Cache.IsOn,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingAuth),
		"url":      "/servers/server/settings/locations/access?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "access",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.Auth != nil && locationConfig.Web.Auth.IsPrior,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingReferers),
		"url":      "/servers/server/settings/locations/referers?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "referer",
		"isOn":     locationConfig.Web != nil && locationConfig.Web.Referers != nil && locationConfig.Web.Referers.IsPrior,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingUserAgents),
		"url":      "/servers/server/settings/locations/userAgent?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "userAgent",
		"isOn":     locationConfig.Web != nil && locationConfig.Web.UserAgent != nil && locationConfig.Web.UserAgent.IsPrior,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingAccessLog),
		"url":      "/servers/server/settings/locations/accessLog?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "accessLog",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.AccessLogRef != nil && locationConfig.Web.AccessLogRef.IsPrior,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingCompression),
		"url":      "/servers/server/settings/locations/compression?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "compression",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.Compression != nil && locationConfig.Web.Compression.IsPrior,
	})

	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingPages),
		"url":      "/servers/server/settings/locations/pages?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "pages",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && (len(locationConfig.Web.Pages) > 0 || (locationConfig.Web.Shutdown != nil && locationConfig.Web.Shutdown.IsPrior)),
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingHTTPHeaders),
		"url":      "/servers/server/settings/locations/headers?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "header",
		"isOn":     locationConfig != nil && this.hasHTTPHeaders(locationConfig.Web),
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingWebsocket),
		"url":      "/servers/server/settings/locations/websocket?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "websocket",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.WebsocketRef != nil && locationConfig.Web.WebsocketRef.IsPrior,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingWebP),
		"url":      "/servers/server/settings/locations/webp?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "webp",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.WebP != nil && locationConfig.Web.WebP.IsPrior,
	})

	menuItems = this.filterMenuItems3(locationConfig, menuItems, serverIdString, locationIdString, secondMenuItem, actionPtr)

	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingStat),
		"url":      "/servers/server/settings/locations/stat?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "stat",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.StatRef != nil && locationConfig.Web.StatRef.IsPrior,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingCharset),
		"url":      "/servers/server/settings/locations/charset?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "charset",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.Charset != nil && locationConfig.Web.Charset.IsPrior,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingRoot),
		"url":      "/servers/server/settings/locations/web?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "web",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.Root != nil && locationConfig.Web.Root.IsPrior,
	})
	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingFastcgi),
		"url":      "/servers/server/settings/locations/fastcgi?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "fastcgi",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.FastcgiRef != nil && locationConfig.Web.FastcgiRef.IsPrior,
	})

	menuItems = this.filterMenuItems2(locationConfig, menuItems, serverIdString, locationIdString, secondMenuItem, actionPtr)

	menuItems = append(menuItems, maps.Map{
		"name":     "-",
		"url":      "",
		"isActive": false,
	})

	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingClientIP),
		"url":      "/servers/server/settings/locations/remoteAddr?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "remoteAddr",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.RemoteAddr != nil && locationConfig.Web.RemoteAddr.IsOn,
	})

	menuItems = append(menuItems, maps.Map{
		"name":     this.Lang(actionPtr, codes.Server_MenuSettingRequestLimit),
		"url":      "/servers/server/settings/locations/requestLimit?serverId=" + serverIdString + "&locationId=" + locationIdString,
		"isActive": secondMenuItem == "requestLimit",
		"isOn":     locationConfig != nil && locationConfig.Web != nil && locationConfig.Web.RequestLimit != nil && locationConfig.Web.RequestLimit.IsOn,
	})

	return menuItems
}

// 检查是否已设置Header
func (this *LocationHelper) hasHTTPHeaders(web *serverconfigs.HTTPWebConfig) bool {
	if web == nil {
		return false
	}
	if web.RequestHeaderPolicyRef != nil && web.RequestHeaderPolicyRef.IsPrior {
		if web.RequestHeaderPolicyRef.IsOn && web.RequestHeaderPolicy != nil && !web.RequestHeaderPolicy.IsEmpty() {
			return true
		}
	}
	if web.ResponseHeaderPolicyRef != nil && web.ResponseHeaderPolicyRef.IsPrior {
		if web.ResponseHeaderPolicyRef.IsOn && web.ResponseHeaderPolicy != nil && !web.ResponseHeaderPolicy.IsEmpty() {
			return true
		}
	}
	return false
}
