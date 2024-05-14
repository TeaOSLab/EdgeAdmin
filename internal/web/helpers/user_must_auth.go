package helpers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/goman"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/index/loginutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/userconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

var nodeLogsCountChanges = make(chan bool, 1)
var ipItemsCountChanges = make(chan bool, 1)

func NotifyNodeLogsCountChange() {
	select {
	case nodeLogsCountChanges <- true:
	default:

	}
}

func NotifyIPItemsCountChanges() {
	select {
	case ipItemsCountChanges <- true:
	default:

	}
}

// 运行日志
var countUnreadNodeLogs int64 = 0
var nodeLogsType = ""

// IP名单
var countUnreadIPItems int64 = 0

// 安全相关
var securityXFFPromptDisabled = false

func init() {
	events.On(events.EventStart, func() {
		// 节点日志数量
		goman.New(func() {
			for range nodeLogsCountChanges {
				rpcClient, err := rpc.SharedRPC()
				if err != nil {
					continue
				}

				countNodeLogsResp, err := rpcClient.NodeLogRPC().CountNodeLogs(rpcClient.Context(0), &pb.CountNodeLogsRequest{
					Role:     nodeconfigs.NodeRoleNode,
					IsUnread: true,
				})
				if err != nil {
					logs.Error(err)
				} else {
					countUnreadNodeLogs = countNodeLogsResp.Count
					if countUnreadNodeLogs > 0 {
						if countUnreadNodeLogs >= 100 {
							countUnreadNodeLogs = 99
						}
						nodeLogsType = "unread"
					}
				}
			}
		})

		// 服务数量
		goman.New(func() {
			for range ipItemsCountChanges {
				rpcClient, err := rpc.SharedRPC()
				if err != nil {
					continue
				}

				countUnreadIPItemsResp, err := rpcClient.IPItemRPC().CountAllEnabledIPItems(rpcClient.Context(0), &pb.CountAllEnabledIPItemsRequest{Unread: true})
				if err != nil {
					logs.Error(err)
				} else {
					countUnreadIPItems = countUnreadIPItemsResp.Count
				}
			}
		})
	})
}

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

	// 检查请求是否合法
	if isEvilRequest(action.Request) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		return false
	}

	// 检测注入
	if !safeFilterRequest(action.Request) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		_, _ = action.ResponseWriter.Write([]byte("Denied By WAF"))
		return false
	}

	// 恢复模式
	if teaconst.IsRecoverMode {
		action.RedirectURL("/recover")
		return false
	}

	// DEMO模式
	if teaconst.IsDemoMode {
		if action.Request.Method == http.MethodPost {
			var actionName = action.Spec.ClassName[strings.LastIndex(action.Spec.ClassName, ".")+1:]
			var denyPrefixes = []string{"Update", "Create", "Delete", "Truncate", "Clean", "Clear", "Reset", "Add", "Remove", "Sync", "Run", "Exec"}
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
	if !checkIP(securityConfig, loginutils.RemoteIP(action)) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		return false
	}

	// 检查请求
	if !checkRequestSecurity(securityConfig, action.Request) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		return false
	}

	// 检查系统是否已经配置过
	if !setup.IsConfigured() {
		action.RedirectURL("/setup")
		return
	}

	var session = action.Session()
	var adminId = session.GetInt64(teaconst.SessionAdminId)

	if adminId <= 0 {
		var errString = session.GetString("@error")
		if len(errString) > 0 {
			action.WriteString("read session failed: " + errString)
			return false
		}
		this.login(action)
		return false
	}

	// 检查指纹
	if securityConfig != nil && securityConfig.CheckClientFingerprint {
		var clientFingerprint = session.GetString("@fingerprint")
		if len(clientFingerprint) > 0 && clientFingerprint != loginutils.CalculateClientFingerprint(action) {
			loginutils.UnsetCookie(action)
			session.Delete()

			this.login(action)
			return false
		}
	}

	// 检查区域
	var oldClientIP = session.GetString("@ip")
	var currentClientIP = loginutils.RemoteIP(action)
	if len(oldClientIP) > 0 && len(currentClientIP) > 0 && oldClientIP != currentClientIP {
		var oldRegion = loginutils.LookupIPRegion(oldClientIP)
		var newRegion = loginutils.LookupIPRegion(currentClientIP)
		if newRegion != oldRegion {
			if securityConfig != nil && securityConfig.CheckClientRegion {
				loginutils.UnsetCookie(action)
				session.Delete()

				this.login(action)
				return false
			} else {
				if !lists.ContainsString([]string{"/messages/badge", "/dns/tasks/check", "/clusters/tasks/check"}, action.Request.URL.Path) {
					// TODO 考虑IP变化时也需要验证，主要是考虑被反向代理的情形
					action.RedirectURL("/login/validate?from=" + url.QueryEscape(action.Request.URL.String()))
					return false
				}
			}
		}
	}

	// 是否正在使用反向代理模式
	if action.Request.Method == http.MethodGet {
		action.Data["teaXFFPrompt"] = false
		if !securityXFFPromptDisabled &&
			(len(action.Header("X-Forwarded-For")) > 0 || len(action.Header("X-Real-Ip")) > 0 || len(action.Header("Cf-Connecting-Ip")) > 0) &&
			securityConfig != nil &&
			len(securityConfig.ClientIPHeaderNames) == 0 {
			action.Data["teaXFFPrompt"] = true
		}
	}

	// 检查用户是否存在
	if !configloaders.CheckAdmin(adminId) {
		loginutils.UnsetCookie(action)
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

	uiConfig, err := configloaders.LoadAdminUIConfig()
	if err != nil {
		action.WriteString(err.Error())
		return false
	}

	// 初始化内置方法
	action.ViewFunc("teaTitle", func() string {
		return action.Data["teaTitle"].(string)
	})

	action.Data["teaShowVersion"] = uiConfig.ShowVersion
	action.Data["teaTitle"] = uiConfig.AdminSystemName
	action.Data["teaName"] = uiConfig.ProductName
	action.Data["teaFaviconFileId"] = uiConfig.FaviconFileId
	action.Data["teaLogoFileId"] = uiConfig.LogoFileId
	action.Data["teaUsername"] = configloaders.FindAdminFullname(adminId)
	action.Data["teaTheme"] = configloaders.FindAdminTheme(adminId)

	action.Data["teaUserAvatar"] = ""

	if !action.Data.Has("teaMenu") {
		action.Data["teaMenu"] = ""
	}

	// 语言
	// Language
	var lang = configloaders.FindAdminLang(adminId)
	if len(lang) == 0 {
		lang = langs.ParseLangFromAction(action)
	}
	action.Data["teaLang"] = lang

	action.Data["teaModules"] = this.modules(lang, actionPtr, adminId, uiConfig)
	action.Data["teaSubMenus"] = []map[string]interface{}{}
	action.Data["teaTabbar"] = []map[string]interface{}{}
	if len(uiConfig.Version) == 0 {
		action.Data["teaVersion"] = teaconst.Version
	} else {
		action.Data["teaVersion"] = uiConfig.Version
	}
	action.Data["teaShowOpenSourceInfo"] = uiConfig.ShowOpenSourceInfo
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
func (this *userMustAuth) modules(langCode string, actionPtr actions.ActionWrapper, adminId int64, adminUIConfig *systemconfigs.AdminUIConfig) []maps.Map {
	// 父级动作
	var action = actionPtr.Object()

	// 未读日志数
	var mainMenu = action.Data.GetString("teaMenu")
	if mainMenu == "clusters" {
		select {
		case nodeLogsCountChanges <- true:
		default:
		}
	} else if mainMenu == "servers" {
		select {
		case ipItemsCountChanges <- true:
		default:
		}
	}

	var result = []maps.Map{}
	for _, m := range FindAllMenuMaps(langCode, nodeLogsType, countUnreadNodeLogs, countUnreadIPItems) {
		if m.GetString("code") == "finance" && !configloaders.ShowFinance() {
			continue
		}

		var module = m.GetString("module")
		if configloaders.AllowModule(adminId, module) {
			if module == "ns" && !adminUIConfig.ContainsModule(userconfigs.UserModuleNS) {
				continue
			}
			if lists.ContainsString([]string{
				configloaders.AdminModuleCodeNode,
				configloaders.AdminModuleCodeDNS,
				configloaders.AdminModuleCodePlan,
				configloaders.AdminModuleCodeServer,
				configloaders.AdminModuleCodeDashboard,
			}, module) && !adminUIConfig.ContainsModule(userconfigs.UserModuleCDN) {
				continue
			}

			result = append(result, m)
		}
	}
	return result
}

// 跳转到登录页
func (this *userMustAuth) login(action *actions.ActionObject) {
	action.RedirectURL("/?from=" + url.QueryEscape(action.Request.RequestURI))
}
