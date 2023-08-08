package recovers

import "github.com/iwind/TeaGo"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(Helper)).
			Prefix("/recover").
			Get("", new(IndexAction)).
			Post("/validateApi", new(ValidateApiAction)).
			Post("/updateHosts", new(UpdateHostsAction)).
			EndAll()
	})
}
