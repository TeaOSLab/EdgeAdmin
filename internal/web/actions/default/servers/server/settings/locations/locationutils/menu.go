// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.
//go:build !plus
// +build !plus

package locationutils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

func filterMenuItems(locationConfig *serverconfigs.HTTPLocationConfig, menuItems []maps.Map, serverIdString string, locationIdString string, secondMenuItem string) []maps.Map {
	return menuItems
}
