package index

import (
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Prefix("/").
			GetPost("", new(IndexAction)).
			EndAll()
	})
}
