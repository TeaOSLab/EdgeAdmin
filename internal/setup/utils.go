package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
)

var isConfigured bool

// IsConfigured 判断系统是否已经配置过
func IsConfigured() bool {
	if isConfigured {
		return true
	}

	_, err := configs.LoadAPIConfig()
	isConfigured = err == nil
	return isConfigured
}

// IsNewInstalled IsNew 检查是否新安装
func IsNewInstalled() bool {
	return configs.IsNewInstalled()
}
