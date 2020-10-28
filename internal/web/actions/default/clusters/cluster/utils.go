package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/iwind/TeaGo/maps"
)

// 安装升级相关的左侧菜单
func LeftMenuItemsForInstall(clusterId int64, selectedItem string) []maps.Map {
	return []maps.Map{
		{
			"name":     "手动安装",
			"url":      "/clusters/cluster/installManual?clusterId=" + numberutils.FormatInt64(clusterId),
			"isActive": selectedItem == "manual",
		},
		{
			"name":     "自动注册",
			"url":      "/clusters/cluster/installNodes?clusterId=" + numberutils.FormatInt64(clusterId),
			"isActive": selectedItem == "register",
		},
		{
			"name":     "远程安装",
			"url":      "/clusters/cluster/installRemote?clusterId=" + numberutils.FormatInt64(clusterId),
			"isActive": selectedItem == "install",
		},
		{
			"name":     "远程升级",
			"url":      "/clusters/cluster/upgradeRemote?clusterId=" + numberutils.FormatInt64(clusterId),
			"isActive": selectedItem == "upgrade",
		},
	}
}
