package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"os"
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(homeDir + "/." + teaconst.ProcessName + "/api.yaml")
	if err != nil {
		return true
	}
	return false
}
