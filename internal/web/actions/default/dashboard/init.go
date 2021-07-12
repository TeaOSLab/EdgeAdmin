package dashboard

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dashboard/boards"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Prefix("/dashboard").
			Data("teaMenu", "dashboard").
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeCommon)).
			GetPost("", new(IndexAction)).

			// 看板
			Prefix("/dashboard/boards").
			Get("", new(boards.IndexAction)).
			Get("/waf", new(boards.WafAction)).
			Post("/wafLogs", new(boards.WafLogsAction)).
			Get("/dns", new(boards.DnsAction)).
			Get("/user", new(boards.UserAction)).

			EndAll()
	})
}
