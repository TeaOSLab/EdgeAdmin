package components

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Data("teaSubMenu", "global").
			Helper(NewHelper()).
			Prefix("/servers/components").
			GetPost("", new(IndexAction)).
			EndAll()
	})
}
