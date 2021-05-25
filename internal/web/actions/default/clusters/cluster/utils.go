package cluster

import (
	"context"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"strconv"
)

// LeftMenuItemsForInstall 安装升级相关的左侧菜单
func LeftMenuItemsForInstall(ctx context.Context, clusterId int64, selectedItem string) []maps.Map {
	rpcClient, _ := rpc.SharedRPC()
	countNotInstalled := int64(0)
	countUpgrade := int64(0)
	if rpcClient != nil {
		{
			resp, err := rpcClient.NodeRPC().CountAllNotInstalledNodesWithNodeClusterId(ctx, &pb.CountAllNotInstalledNodesWithNodeClusterIdRequest{NodeClusterId: clusterId})
			if err == nil {
				countNotInstalled = resp.Count
			}
		}
		{
			resp, err := rpcClient.NodeRPC().CountAllUpgradeNodesWithNodeClusterId(ctx, &pb.CountAllUpgradeNodesWithNodeClusterIdRequest{NodeClusterId: clusterId})
			if err == nil {
				countUpgrade = resp.Count
			}
		}
	}

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
			"name":     "远程安装(" + strconv.FormatInt(countNotInstalled, 10) + ")",
			"url":      "/clusters/cluster/installRemote?clusterId=" + numberutils.FormatInt64(clusterId),
			"isActive": selectedItem == "install",
		},
		{
			"name":     "远程升级(" + strconv.FormatInt(countUpgrade, 10) + ")",
			"url":      "/clusters/cluster/upgradeRemote?clusterId=" + numberutils.FormatInt64(clusterId),
			"isActive": selectedItem == "upgrade",
		},
	}
}
