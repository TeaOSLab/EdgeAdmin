package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeAdmin)).
			Data("teaMenu", "admins").
			Data("teaSubMenu", "recipients").
			Prefix("/admins/recipients/tasks").
			Get("", new(IndexAction)).
			Post("/taskInfo", new(TaskInfoAction)).
			Post("/delete", new(DeleteAction)).
			EndAll()
	})
}
