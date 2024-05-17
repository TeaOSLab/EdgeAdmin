// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.
//go:build !plus

package serverutils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

func (this *ServerHelper) filterMenuItems(serverConfig *serverconfigs.ServerConfig, menuItems []maps.Map, serverIdString string, secondMenuItem string, actionPtr actions.ActionWrapper) []maps.Map {
	return menuItems
}

func (this *ServerHelper) filterMenuItems2(serverConfig *serverconfigs.ServerConfig, menuItems []maps.Map, serverIdString string, secondMenuItem string, actionPtr actions.ActionWrapper) []maps.Map {
	return menuItems
}

func (this *ServerHelper) filterMenuItems3(serverConfig *serverconfigs.ServerConfig, menuItems []maps.Map, serverIdString string, secondMenuItem string, actionPtr actions.ActionWrapper) []maps.Map {
	return menuItems
}
