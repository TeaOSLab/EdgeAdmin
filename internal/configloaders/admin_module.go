package configloaders

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/maps"
)

type AdminModuleCode = string

const (
	AdminModuleCodeServer  AdminModuleCode = "server"  // 网站
	AdminModuleCodeNode    AdminModuleCode = "node"    // 节点
	AdminModuleCodeDNS     AdminModuleCode = "dns"     // DNS
	AdminModuleCodeAdmin   AdminModuleCode = "admin"   // 系统用户
	AdminModuleCodeUser    AdminModuleCode = "user"    // 平台用户
	AdminModuleCodeFinance AdminModuleCode = "finance" // 财务
	AdminModuleCodeLog     AdminModuleCode = "log"     // 日志
	AdminModuleCodeSetting AdminModuleCode = "setting" // 设置
	AdminModuleCodeCommon  AdminModuleCode = "common"  // 只要登录就可以访问的模块
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

// 检查用户是否存在
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

// 检查模块是否允许访问
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

// 获取管理员第一个可访问模块
func FindFirstAdminModule(adminId int64) (module AdminModuleCode, ok bool) {
	locker.Lock()
	defer locker.Unlock()
	list, ok2 := sharedAdminModuleMapping[adminId]
	if ok2 {
		if list.IsSuper {
			return AdminModuleCodeServer, true
		} else if len(list.Modules) > 0 {
			return list.Modules[0].Code, true
		}
	}
	return
}

// 查找某个管理员名称
func FindAdminFullname(adminId int64) string {
	locker.Lock()
	defer locker.Unlock()

	list, ok := sharedAdminModuleMapping[adminId]
	if ok {
		return list.Fullname
	}
	return ""
}

// 所有权限列表
func AllModuleMaps() []maps.Map {
	return []maps.Map{
		{
			"name": "网站服务",
			"code": AdminModuleCodeServer,
			"url":  "/servers",
		},
		{
			"name": "边缘节点",
			"code": AdminModuleCodeNode,
			"url":  "/clusters",
		},
		{
			"name": "域名解析",
			"code": AdminModuleCodeDNS,
			"url":  "/dns",
		},
		{
			"name": "平台用户",
			"code": AdminModuleCodeUser,
			"url":  "/users",
		},
		{
			"name": "系统用户",
			"code": AdminModuleCodeAdmin,
			"url":  "/admins",
		},
		{
			"name": "财务管理",
			"code": AdminModuleCodeFinance,
			"url":  "/finance",
		},
		{
			"name": "日志审计",
			"code": AdminModuleCodeLog,
			"url":  "/log",
		},
		{
			"name": "系统设置",
			"code": AdminModuleCodeSetting,
			"url":  "/settings",
		},
	}
}
