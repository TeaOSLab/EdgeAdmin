package logs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeNS)).
			Data("teaMenu", "ns").
			Data("teaSubMenu", "accessLog").
			Prefix("/ns/clusters/accessLogs").
			Get("", new(IndexAction)).
			EndAll()
	})
}
