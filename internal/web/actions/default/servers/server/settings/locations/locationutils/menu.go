// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.
//go:build !plus

package locationutils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

func (this *LocationHelper) filterMenuItems1(locationConfig *serverconfigs.HTTPLocationConfig, menuItems []maps.Map, serverIdString string, locationIdString string, secondMenuItem string, actionPtr actions.ActionWrapper) []maps.Map {
	return menuItems
}

func (this *LocationHelper) filterMenuItems2(locationConfig *serverconfigs.HTTPLocationConfig, menuItems []maps.Map, serverIdString string, locationIdString string, secondMenuItem string, actionPtr actions.ActionWrapper) []maps.Map {
	return menuItems
}

func (this *LocationHelper) filterMenuItems3(locationConfig *serverconfigs.HTTPLocationConfig, menuItems []maps.Map, serverIdString string, locationIdString string, secondMenuItem string, actionPtr actions.ActionWrapper) []maps.Map {
	return menuItems
}
