// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type ThemeAction struct {
	actionutils.ParentAction
}

func (this *ThemeAction) RunPost(params struct{}) {
	theme := configloaders.FindAdminTheme(this.AdminId())

	var themes = []string{"theme1", "theme2", "theme3", "theme4", "theme5", "theme6", "theme7"}
	var nextTheme = "theme1"
	if len(theme) == 0 {
		nextTheme = "theme2"
	} else {
		for index, t := range themes {
			if t == theme {
				if index < len(themes)-1 {
					nextTheme = themes[index+1]
					break
				}
			}
		}
	}

	_, err := this.RPC().AdminRPC().UpdateAdminTheme(this.AdminContext(), &pb.UpdateAdminThemeRequest{
		AdminId: this.AdminId(),
		Theme:   nextTheme,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	configloaders.UpdateAdminTheme(this.AdminId(), nextTheme)

	this.Data["theme"] = nextTheme

	this.Success()
}
