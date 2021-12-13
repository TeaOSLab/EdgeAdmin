package helpers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net"
	"net/http"
	"net/url"
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
	action.Data["teaModules"] = this.modules(actionPtr, adminId)
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
func (this *userMustAuth) modules(actionPtr actions.ActionWrapper, adminId int64) []maps.Map {
	var countUnreadNodeLogs int64 = 0
	var nodeLogsType = ""

	// 父级动作
	parentAction, ok := actionPtr.(actionutils.ActionInterface)
	if ok {
		var action = actionPtr.Object()

		// 未读日志数
		if action.Data.GetString("teaMenu") == "clusters" {
			countNodeLogsResp, err := parentAction.RPC().NodeLogRPC().CountNodeLogs(parentAction.AdminContext(), &pb.CountNodeLogsRequest{
				Role:     nodeconfigs.NodeRoleNode,
				IsUnread: true,
			})
			if err != nil {
				logs.Error(err)
			} else {
				var countNodeLogs = countNodeLogsResp.Count
				if countNodeLogs > 0 {
					countUnreadNodeLogs = countNodeLogs
					if countUnreadNodeLogs >= 100 {
						countUnreadNodeLogs = 99
					}
					nodeLogsType = "unread"
				}
			}
		}
	}

	result := []maps.Map{}
	for _, m := range FindAllMenuMaps(nodeLogsType, countUnreadNodeLogs) {
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
	action.RedirectURL("/?from=" + url.QueryEscape(action.Request.RequestURI))
}
