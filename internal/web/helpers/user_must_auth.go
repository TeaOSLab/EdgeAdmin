package helpers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net"
	"net/http"
	"reflect"
	"strings"
)

// 认证拦截
type userMustAuth struct {
	AdminId int64
	module  string
}

func NewUserMustAuth(module string) *userMustAuth {
	return &userMustAuth{module: module}
}

func (this *userMustAuth) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	var action = actionPtr.Object()

	// 恢复模式
	if teaconst.IsRecoverMode {
		action.RedirectURL("/recover")
		return false
	}

	// DEMO模式
	if teaconst.IsDemoMode {
		if action.Request.Method == http.MethodPost {
			var actionName = action.Spec.ClassName[strings.LastIndex(action.Spec.ClassName, ".")+1:]
			var denyPrefixes = []string{"Update", "Create", "Delete", "Truncate", "Clean", "Clear", "Reset", "Add", "Remove", "Sync"}
			for _, prefix := range denyPrefixes {
				if strings.HasPrefix(actionName, prefix) {
					action.Fail(teaconst.ErrorDemoOperation)
					return false
				}
			}

			if strings.Index(action.Spec.PkgPath, "settings") > 0 || strings.Index(action.Spec.PkgPath, "delete") > 0 || strings.Index(action.Spec.PkgPath, "update") > 0 {
				action.Fail(teaconst.ErrorDemoOperation)
				return false
			}
		}
	}

	// 安全相关
	securityConfig, _ := configloaders.LoadSecurityConfig()
	if securityConfig == nil {
		action.AddHeader("X-Frame-Options", "SAMEORIGIN")
	} else if len(securityConfig.Frame) > 0 {
		action.AddHeader("X-Frame-Options", securityConfig.Frame)
	}
	action.AddHeader("Content-Security-Policy", "default-src 'self' data:; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'")

	// 检查IP
	if !checkIP(securityConfig, action.RequestRemoteIP()) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		return false
	}
	remoteAddr, _, _ := net.SplitHostPort(action.Request.RemoteAddr)
	if len(remoteAddr) > 0 && remoteAddr != action.RequestRemoteIP() && !checkIP(securityConfig, remoteAddr) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		return false
	}

	// 检查系统是否已经配置过
	if !setup.IsConfigured() {
		action.RedirectURL("/setup")
		return
	}

	var session = action.Session()
	var adminId = session.GetInt64("adminId")

	if adminId <= 0 {
		this.login(action)
		return false
	}

	// 检查用户是否存在
	if !configloaders.CheckAdmin(adminId) {
		session.Delete()

		this.login(action)
		return false
	}

	// 检查用户权限
	if len(this.module) > 0 && !configloaders.AllowModule(adminId, this.module) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		action.WriteString("Permission Denied.")
		return false
	}

	this.AdminId = adminId
	action.Context.Set("adminId", this.AdminId)

	if action.Request.Method != http.MethodGet {
		return true
	}

	config, err := configloaders.LoadAdminUIConfig()
	if err != nil {
		action.WriteString(err.Error())
		return false
	}

	// 初始化内置方法
	action.ViewFunc("teaTitle", func() string {
		return action.Data["teaTitle"].(string)
	})

	action.Data["teaShowVersion"] = config.ShowVersion
	action.Data["teaTitle"] = config.AdminSystemName
	action.Data["teaName"] = config.ProductName
	action.Data["teaFaviconFileId"] = config.FaviconFileId
	action.Data["teaLogoFileId"] = config.LogoFileId
	action.Data["teaUsername"] = configloaders.FindAdminFullname(adminId)
	action.Data["teaTheme"] = configloaders.FindAdminTheme(adminId)

	action.Data["teaUserAvatar"] = ""

	if !action.Data.Has("teaMenu") {
		action.Data["teaMenu"] = ""
	}
	action.Data["teaModules"] = this.modules(adminId)
	action.Data["teaSubMenus"] = []map[string]interface{}{}
	action.Data["teaTabbar"] = []map[string]interface{}{}
	if len(config.Version) == 0 {
		action.Data["teaVersion"] = teaconst.Version
	} else {
		action.Data["teaVersion"] = config.Version
	}
	action.Data["teaShowOpenSourceInfo"] = config.ShowOpenSourceInfo
	action.Data["teaIsSuper"] = false
	action.Data["teaIsPlus"] = teaconst.IsPlus
	action.Data["teaDemoEnabled"] = teaconst.IsDemoMode
	action.Data["teaShowFinance"] = configloaders.ShowFinance()
	if !action.Data.Has("teaSubMenu") {
		action.Data["teaSubMenu"] = ""
	}
	action.Data["teaCheckNodeTasks"] = configloaders.AllowModule(adminId, configloaders.AdminModuleCodeNode)
	action.Data["teaCheckDNSTasks"] = configloaders.AllowModule(adminId, configloaders.AdminModuleCodeDNS)

	// 菜单
	action.Data["firstMenuItem"] = ""

	// 未读消息数
	action.Data["teaBadge"] = 0

	// 调用Init
	initMethod := reflect.ValueOf(actionPtr).MethodByName("Init")
	if initMethod.IsValid() {
		initMethod.Call([]reflect.Value{})
	}

	return true
}

// 菜单配置
func (this *userMustAuth) modules(adminId int64) []maps.Map {
	allMaps := []maps.Map{
		{
			"code":   "dashboard",
			"module": configloaders.AdminModuleCodeDashboard,
			"name":   "数据看板",
			"icon":   "dashboard",
		},
		{
			"code":     "servers",
			"module":   configloaders.AdminModuleCodeServer,
			"name":     "网站服务",
			"subtitle": "服务列表",
			"icon":     "clone outsize",
			"subItems": []maps.Map{
				{
					"name": "服务分组",
					"url":  "/servers/groups",
					"code": "group",
				},
				{
					"name": "证书管理",
					"url":  "/servers/certs",
					"code": "cert",
				},
				{
					"name": "访问日志",
					"url":  "/servers/logs",
					"code": "log",
				},
				{
					"name": "缓存策略",
					"url":  "/servers/components/cache",
					"code": "cache",
				},
				{
					"name": "WAF策略",
					"url":  "/servers/components/waf",
					"code": "waf",
				},
				{
					"name": "日志策略",
					"url":  "/servers/accesslogs",
					"code": "accesslog",
					"isOn": teaconst.IsPlus,
				},
				{
					"name": "IP名单",
					"url":  "/servers/iplists",
					"code": "iplist",
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
					"name": "运行日志",
					"url":  "/clusters/logs",
					"code": "log",
				},
				{
					"name": "IP地址",
					"url":  "/clusters/ip-addrs",
					"code": "ipAddr",
					"isOn": teaconst.IsPlus,
				},
				{
					"name": "区域监控",
					"url":  "/clusters/monitors",
					"code": "monitor",
					"isOn": teaconst.IsPlus,
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
					"name": "问题修复",
					"url":  "/dns/issues",
					"code": "issue",
				},
				{
					"name": "DNS服务商",
					"url":  "/dns/providers",
					"code": "provider",
				},
			},
		},
		{
			"code":   "ns",
			"module": configloaders.AdminModuleCodeNS,
			"name":   "智能DNS",
			"icon":   "cubes",
			"isOn":   teaconst.IsPlus,
			"subItems": []maps.Map{
				{
					"name": "域名管理",
					"url":  "/ns/domains",
					"code": "domain",
				},
				{
					"name": "集群管理",
					"url":  "/ns/clusters",
					"code": "cluster",
				},
				{
					"name": "线路管理",
					"url":  "/ns/routes",
					"code": "route",
				},
				{
					"name": "访问日志",
					"url":  "/ns/clusters/accessLogs",
					"code": "accessLog",
				},
				{
					"name": "运行日志",
					"url":  "/ns/clusters/logs",
					"code": "log",
				},
				{
					"name": "全局配置",
					"url":  "/ns/settings",
					"code": "setting",
				},
				{
					"name": "解析测试",
					"url":  "/ns/test",
					"code": "test",
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
			"code":   "finance",
			"module": configloaders.AdminModuleCodeFinance,
			"name":   "财务管理",
			"icon":   "yen sign",
			"isOn":   teaconst.IsPlus,
		},
		{
			"code":     "admins",
			"module":   configloaders.AdminModuleCodeAdmin,
			"name":     "系统用户",
			"subtitle": "用户列表",
			"icon":     "user secret",
			"subItems": []maps.Map{
				{
					"name": "通知媒介",
					"url":  "/admins/recipients",
					"code": "recipients",
					"isOn": teaconst.IsPlus,
				},
			},
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

	result := []maps.Map{}
	for _, m := range allMaps {
		if m.GetString("code") == "finance" && !configloaders.ShowFinance() {
			continue
		}

		module := m.GetString("module")
		if configloaders.AllowModule(adminId, module) {
			result = append(result, m)
		}
	}
	return result
}

// 跳转到登录页
func (this *userMustAuth) login(action *actions.ActionObject) {
	action.RedirectURL("/")
}
