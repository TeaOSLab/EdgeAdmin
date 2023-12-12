package configloaders

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AdminModuleCode = string

const (
	AdminModuleCodeDashboard AdminModuleCode = "dashboard" // 看板
	AdminModuleCodeServer    AdminModuleCode = "server"    // 网站
	AdminModuleCodeNode      AdminModuleCode = "node"      // 节点
	AdminModuleCodeDNS       AdminModuleCode = "dns"       // DNS
	AdminModuleCodeNS        AdminModuleCode = "ns"        // 域名服务
	AdminModuleCodeAdmin     AdminModuleCode = "admin"     // 系统用户
	AdminModuleCodeUser      AdminModuleCode = "user"      // 平台用户
	AdminModuleCodeFinance   AdminModuleCode = "finance"   // 财务
	AdminModuleCodePlan      AdminModuleCode = "plan"      // 套餐
	AdminModuleCodeLog       AdminModuleCode = "log"       // 日志
	AdminModuleCodeSetting   AdminModuleCode = "setting"   // 设置
	AdminModuleCodeTicket    AdminModuleCode = "ticket"    // 工单
	AdminModuleCodeCommon    AdminModuleCode = "common"    // 只要登录就可以访问的模块
)

var sharedAdminModuleMapping = map[int64]*AdminModuleList{} // adminId => AdminModuleList

func loadAdminModuleMapping() (map[int64]*AdminModuleList, error) {
	if len(sharedAdminModuleMapping) > 0 {
		return sharedAdminModuleMapping, nil
	}

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	modulesResp, err := rpcClient.AdminRPC().FindAllAdminModules(rpcClient.Context(0), &pb.FindAllAdminModulesRequest{})
	if err != nil {
		return nil, err
	}
	mapping := map[int64]*AdminModuleList{}
	for _, m := range modulesResp.AdminModules {
		list := &AdminModuleList{
			IsSuper:  m.IsSuper,
			Fullname: m.Fullname,
			Theme:    m.Theme,
			Lang:     m.Lang,
		}

		for _, pbModule := range m.Modules {
			list.Modules = append(list.Modules, &systemconfigs.AdminModule{
				Code:     pbModule.Code,
				AllowAll: pbModule.AllowAll,
				Actions:  pbModule.Actions,
			})
		}

		mapping[m.AdminId] = list
	}

	sharedAdminModuleMapping = mapping

	return sharedAdminModuleMapping, nil
}

func NotifyAdminModuleMappingChange() error {
	locker.Lock()
	defer locker.Unlock()
	sharedAdminModuleMapping = map[int64]*AdminModuleList{}
	_, err := loadAdminModuleMapping()
	return err
}

// CheckAdmin 检查用户是否存在
func CheckAdmin(adminId int64) bool {
	locker.Lock()
	defer locker.Unlock()

	// 如果还没有数据，则尝试加载
	if len(sharedAdminModuleMapping) == 0 {
		_, _ = loadAdminModuleMapping()
	}

	_, ok := sharedAdminModuleMapping[adminId]
	return ok
}

// AllowModule 检查模块是否允许访问
func AllowModule(adminId int64, module string) bool {
	locker.Lock()
	defer locker.Unlock()

	if module == AdminModuleCodeCommon {
		return true
	}

	if len(sharedAdminModuleMapping) == 0 {
		_, _ = loadAdminModuleMapping()
	}

	list, ok := sharedAdminModuleMapping[adminId]
	if ok {
		return list.Allow(module)
	}

	return false
}

// FindFirstAdminModule 获取管理员第一个可访问模块
func FindFirstAdminModule(adminId int64) (module AdminModuleCode, ok bool) {
	locker.Lock()
	defer locker.Unlock()
	list, ok2 := sharedAdminModuleMapping[adminId]
	if ok2 {
		if list.IsSuper {
			return AdminModuleCodeDashboard, true
		} else if len(list.Modules) > 0 {
			return list.Modules[0].Code, true
		}
	}
	return
}

// FindAdminFullname 查找某个管理员名称
func FindAdminFullname(adminId int64) string {
	locker.Lock()
	defer locker.Unlock()

	list, ok := sharedAdminModuleMapping[adminId]
	if ok {
		return list.Fullname
	}
	return ""
}

// FindAdminTheme 查找某个管理员选择的风格
func FindAdminTheme(adminId int64) string {
	locker.Lock()
	defer locker.Unlock()

	list, ok := sharedAdminModuleMapping[adminId]
	if ok {
		return list.Theme
	}
	return ""
}

// UpdateAdminTheme 设置某个管理员的风格
func UpdateAdminTheme(adminId int64, theme string) {
	locker.Lock()
	defer locker.Unlock()

	list, ok := sharedAdminModuleMapping[adminId]
	if ok {
		list.Theme = theme
	}
}

// FindAdminLang 查找某个管理员选择的语言
func FindAdminLang(adminId int64) string {
	locker.Lock()
	defer locker.Unlock()

	list, ok := sharedAdminModuleMapping[adminId]
	if ok {
		return list.Lang
	}
	return ""
}

// UpdateAdminLang 修改某个管理员选择的语言
func UpdateAdminLang(adminId int64, langCode string) {
	locker.Lock()
	defer locker.Unlock()

	list, ok := sharedAdminModuleMapping[adminId]
	if ok {
		list.Lang = langCode
	}
}

func FindAdminLangForAction(actionPtr actions.ActionWrapper) (langCode langs.LangCode) {
	locker.Lock()
	defer locker.Unlock()

	var adminId = actionPtr.Object().Session().GetInt64(teaconst.SessionAdminId)
	list, ok := sharedAdminModuleMapping[adminId]
	var result = ""
	if ok {
		result = list.Lang
	}

	if len(result) == 0 {
		result = langs.ParseLangFromAction(actionPtr)
	}

	return result
}

// AllModuleMaps 所有权限列表
func AllModuleMaps(langCode string) []maps.Map {
	var m = []maps.Map{
		{
			"name": langs.Message(langCode, codes.AdminMenu_Dashboard),
			"code": AdminModuleCodeDashboard,
			"url":  "/dashboard",
		},
		{
			"name": langs.Message(langCode, codes.AdminMenu_Servers),
			"code": AdminModuleCodeServer,
			"url":  "/servers",
		},
		{
			"name": langs.Message(langCode, codes.AdminMenu_Nodes),
			"code": AdminModuleCodeNode,
			"url":  "/clusters",
		},
		{
			"name": langs.Message(langCode, codes.AdminMenu_DNS),
			"code": AdminModuleCodeDNS,
			"url":  "/dns",
		},
	}
	if teaconst.IsPlus {
		m = append(m, maps.Map{
			"name": langs.Message(langCode, codes.AdminMenu_NS),
			"code": AdminModuleCodeNS,
			"url":  "/ns",
		})
	}
	m = append(m, []maps.Map{
		{
			"name": langs.Message(langCode, codes.AdminMenu_Users),
			"code": AdminModuleCodeUser,
			"url":  "/users",
		},
		{
			"name": langs.Message(langCode, codes.AdminMenu_Admins),
			"code": AdminModuleCodeAdmin,
			"url":  "/admins",
		},
		{
			"name": langs.Message(langCode, codes.AdminMenu_Finance),
			"code": AdminModuleCodeFinance,
			"url":  "/finance",
		},
	}...)

	if teaconst.IsPlus {
		m = append(m, []maps.Map{
			{
				"name": langs.Message(langCode, codes.AdminMenu_Plans),
				"code": AdminModuleCodePlan,
				"url":  "/plans",
			},
			{
				"name": langs.Message(langCode, codes.AdminMenu_Tickets),
				"code": AdminModuleCodeTicket,
				"url":  "/tickets",
			},
		}...)
	}

	m = append(m, []maps.Map{
		{
			"name": langs.Message(langCode, codes.AdminMenu_Logs),
			"code": AdminModuleCodeLog,
			"url":  "/log",
		},
		{
			"name": langs.Message(langCode, codes.AdminMenu_Settings),
			"code": AdminModuleCodeSetting,
			"url":  "/settings",
		},
	}...)
	return m
}
