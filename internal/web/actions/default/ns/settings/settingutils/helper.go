// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package settingutils

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type Helper struct {
}

func (this *Helper) BeforeAction(actionPtr actions.ActionWrapper) (goNext bool) {
	var action = actionPtr.Object()
	secondMenuItem := action.Data.GetString("secondMenuItem")
	action.Data["leftMenuItems"] = this.createSettingMenu(secondMenuItem)
	return true
}

func (this *Helper) createSettingMenu(selectedItem string) (items []maps.Map) {
	return []maps.Map{
		{
			"name":     "访问日志",
			"url":      "/ns/settings/accesslogs",
			"isActive": selectedItem == "accessLog",
		},
	}
}
