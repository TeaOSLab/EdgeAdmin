package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeServer)).
			Data("teaMenu", "servers").
			Data("teaSubMenu", "iplist").
			Prefix("/servers/iplists").
			Get("", new(IndexAction)).
			Get("/lists", new(ListsAction)).
			GetPost("/createPopup", new(CreatePopupAction)).
			Get("/list", new(ListAction)).
			GetPost("/import", new(ImportAction)).
			GetPost("/export", new(ExportAction)).
			Get("/exportData", new(ExportDataAction)).
			Post("/delete", new(DeleteAction)).
			Post("/deleteItems", new(DeleteItemsAction)).
			Post("/deleteCount", new(DeleteCountAction)).
			GetPost("/test", new(TestAction)).
			GetPost("/update", new(UpdateAction)).
			Get("/items", new(ItemsAction)).
			Get("/selectPopup", new(SelectPopupAction)).

			// IP相关
			GetPost("/createIPPopup", new(CreateIPPopupAction)).
			GetPost("/updateIPPopup", new(UpdateIPPopupAction)).
			Post("/deleteIP", new(DeleteIPAction)).
			Get("/accessLogsPopup", new(AccessLogsPopupAction)).
			Post("/readAll", new(ReadAllAction)).

			// 防火墙
			GetPost("/bindHTTPFirewallPopup", new(BindHTTPFirewallPopupAction)).
			Post("/unbindHTTPFirewall", new(UnbindHTTPFirewallAction)).
			Post("/httpFirewall", new(HttpFirewallAction)).

			// 选项数据
			Post("/levelOptions", new(LevelOptionsAction)).
			EndAll()
	})
}
