package bills

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeFinance)).
			Data("teaMenu", "finance").

			// 财务管理
			Prefix("/finance/bills").
			Get("", new(IndexAction)).
			GetPost("/generate", new(GenerateAction)).

			EndAll()
	})
}
