// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.
//go:build !plus

package serverutils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

func filterMenuItems(serverConfig *serverconfigs.ServerConfig, menuItems []maps.Map, serverIdString string, secondMenuItem string) []maps.Map {
	return menuItems
}
