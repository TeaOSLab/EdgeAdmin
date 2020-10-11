package ui

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/actions"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/ui").
			Get("/download", new(DownloadAction)).

			// 以下的需要压缩
			Helper(new(actions.Gzip)).
			Get("/components.js", new(ComponentsAction)).
			EndAll()
	})
}
