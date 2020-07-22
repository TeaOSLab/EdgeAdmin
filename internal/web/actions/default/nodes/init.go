package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Prefix("/nodes").
			Get("", new(IndexAction)).
			GetPost("/create", new(CreateAction)).
			EndAll()
	})
}
