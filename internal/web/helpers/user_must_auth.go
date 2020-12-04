package helpers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	nodes "github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"reflect"
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

	// 检查用户权限
	if len(this.module) > 0 && !configloaders.AllowModule(adminId, this.module) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		action.WriteString("Permission Denied.")
		return false
	}

	// 检查用户是否存在
	rpc, err := nodes.SharedRPC()
	if err != nil {
		action.WriteString("setup rpc error: " + err.Error())
		utils.PrintError(err)
		return false
	}

	rpcResp, err := rpc.AdminRPC().CheckAdminExists(rpc.Context(0), &pb.CheckAdminExistsRequest{AdminId: adminId})
	if err != nil {
		utils.PrintError(err)
		action.WriteString(teaconst.ErrServer)
		return false
	}

	if !rpcResp.IsOk {
		session.Delete()

		this.login(action)
		return false
	}

	this.AdminId = adminId
	action.Context.Set("adminId", this.AdminId)

	if action.Request.Method != http.MethodGet {
		return true
	}

	config, err := configloaders.LoadUIConfig()
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

	resp, err := rpc.AdminRPC().FindAdminFullname(rpc.Context(0), &pb.FindAdminFullnameRequest{AdminId: this.AdminId})
	if err != nil {
		utils.PrintError(err)
		action.Data["teaUsername"] = ""
	} else {
		action.Data["teaUsername"] = resp.Fullname
	}

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
	action.Data["teaDemoEnabled"] = teaconst.IsDemo
	if !action.Data.Has("teaSubMenu") {
		action.Data["teaSubMenu"] = ""
	}

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
			"code":   "servers",
			"module": configloaders.AdminModuleCodeServer,
			"name":   "网站服务",
			"icon":   "clone outsize",
			"subItems": []maps.Map{
				{
					"name": "通用设置",
					"url":  "/servers/components",
					"code": "global",
				},
				{
					"name": "服务分组",
					"url":  "/servers/components/groups",
					"code": "group",
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
					"name": "证书管理",
					"url":  "/servers/certs",
					"code": "cert",
				},
			},
		},
		{
			"code":   "clusters",
			"module": configloaders.AdminModuleCodeNode,
			"name":   "边缘节点",
			"icon":   "cloud",
			"subItems": []maps.Map{
				{
					"name": "SSH认证",
					"url":  "/clusters/grants",
					"code": "grant",
				},
			},
		},
		{
			"code":   "dns",
			"module": configloaders.AdminModuleCodeDNS,
			"name":   "域名解析",
			"icon":   "globe",
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
			"code":   "users",
			"module": configloaders.AdminModuleCodeUser,
			"name":   "平台用户",
			"icon":   "users",
		},
		{
			"code":   "admins",
			"module": configloaders.AdminModuleCodeAdmin,
			"name":   "系统用户",
			"icon":   "user secret",
		},
		{
			"code":   "log",
			"module": configloaders.AdminModuleCodeLog,
			"name":   "日志审计",
			"icon":   "history",
		},
		{
			"code":   "settings",
			"module": configloaders.AdminModuleCodeSetting,
			"name":   "系统设置",
			"icon":   "setting",
		},
	}

	result := []maps.Map{}
	for _, m := range allMaps {
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
