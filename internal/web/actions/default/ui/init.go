package ui

import (
	"compress/gzip"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/actions"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/ui").
			Get("/download", new(DownloadAction)).
			GetPost("/selectProvincesPopup", new(SelectProvincesPopupAction)).
			GetPost("/selectCountriesPopup", new(SelectCountriesPopupAction)).

			// 以下的需要压缩
			Helper(&actions.Gzip{Level: gzip.BestCompression}).
			Get("/components.js", new(ComponentsAction)).
			EndAll()
	})
}
