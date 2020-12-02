package configloaders

import "github.com/iwind/TeaGo/maps"

type AdminModuleCode = string

const (
	AdminModuleCodeServer  AdminModuleCode = "server"
	AdminModuleCodeNode    AdminModuleCode = "node"
	AdminModuleCodeDNS     AdminModuleCode = "dns"
	AdminModuleCodeAdmin   AdminModuleCode = "admin"
	AdminModuleCodeLog     AdminModuleCode = "log"
	AdminModuleCodeSetting AdminModuleCode = "setting"
)

var adminModuleMapping = map[int64]*AdminModuleList{} // adminId => AdminModuleList

func LoadAdminModuleMapping() (map[int64]*AdminModuleList, error) {
	locker.Lock()
	defer locker.Unlock()

	if len(adminModuleMapping) > 0 {
		return adminModuleMapping, nil
	}

	// TODO

	return nil, nil
}

func NotifyAdminModuleMappingChange() error {
	locker.Lock()
	adminModuleMapping = map[int64]*AdminModuleList{}
	locker.Unlock() // 这里结束是为了避免和LoadAdminModuleMapping()造成死锁
	_, err := LoadAdminModuleMapping()
	return err
}

func IsAllowModule(adminId int64, module string) bool {
	// TODO
	return false
}

// 所有权限列表
func AllModuleMaps() []maps.Map {
	return []maps.Map{
		{
			"name": "网站服务",
			"code": AdminModuleCodeServer,
		},
		{
			"name": "边缘节点",
			"code": AdminModuleCodeNode,
		},
		{
			"name": "域名解析",
			"code": AdminModuleCodeDNS,
		},
		{
			"name": "系统用户",
			"code": AdminModuleCodeAdmin,
		},
		{
			"name": "日志审计",
			"code": AdminModuleCodeLog,
		},
		{
			"name": "系统设置",
			"code": AdminModuleCodeSetting,
		},
	}
}
