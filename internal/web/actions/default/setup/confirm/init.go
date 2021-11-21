package confirm

import "github.com/iwind/TeaGo"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(Helper)).
			Prefix("/setup/confirm").
			GetPost("", new(IndexAction)).
			EndAll()
	})
}
