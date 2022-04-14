// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

//go:build !plus
// +build !plus

package servergrouputils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

func filterMenuItems(leftMenuItems []maps.Map, groupId int64, urlPrefix string, menuItem string, configInfoResp *pb.FindEnabledServerGroupConfigInfoResponse) []maps.Map {
	return leftMenuItems
}

func filterMenuItems2(leftMenuItems []maps.Map, groupId int64, urlPrefix string, menuItem string, configInfoResp *pb.FindEnabledServerGroupConfigInfoResponse) []maps.Map {
	return leftMenuItems
}
