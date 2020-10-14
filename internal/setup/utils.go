package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
)

var isConfigured bool

// 判断系统是否已经配置过
// TODO 检查节点版本和数据库版本是否一致，如果不一致则跳转到升级页面
func IsConfigured() bool {
	if isConfigured {
		return true
	}

	_, err := configs.LoadAPIConfig()
	isConfigured = err == nil
	return isConfigured
}
