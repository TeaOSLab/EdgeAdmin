package about

import "github.com/iwind/TeaGo"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/about").
			Get("/qq", new(QqAction)).
			EndAll()
	})
}
