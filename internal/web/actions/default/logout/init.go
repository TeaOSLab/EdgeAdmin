package logout

import "github.com/iwind/TeaGo"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/logout").
			Get("", new(IndexAction)).
			EndAll()
	})
}
