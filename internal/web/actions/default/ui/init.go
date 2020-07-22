package ui

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/actions"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(actions.Gzip)).
			Prefix("/ui").
			Get("/components.js", new(ComponentsAction)).
			EndAll()
	})
}
