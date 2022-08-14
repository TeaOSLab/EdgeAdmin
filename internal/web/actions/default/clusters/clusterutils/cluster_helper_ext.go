// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package clusterutils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

func filterMenuItems1(items []maps.Map, info *pb.FindEnabledNodeClusterConfigInfoResponse, clusterIdString string, selectedItem string) []maps.Map {
	return items
}

func filterMenuItems2(items []maps.Map, info *pb.FindEnabledNodeClusterConfigInfoResponse, clusterIdString string, selectedItem string) []maps.Map {
	return items
}
