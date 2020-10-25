package csrf

import "github.com/iwind/TeaGo"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/csrf").
			Get("/token", new(TokenAction)).
			EndAll()
	})
}
