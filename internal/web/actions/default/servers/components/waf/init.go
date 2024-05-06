package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/waf/ipadmin"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeServer)).
			Helper(NewHelper()).
			Data("teaMenu", "servers").
			Data("teaSubMenu", "waf").
			Prefix("/servers/components/waf").
			Get("", new(IndexAction)).
			GetPost("/createPopup", new(CreatePopupAction)).
			Post("/delete", new(DeleteAction)).
			Get("/policy", new(PolicyAction)).
			Post("/upgradeTemplate", new(UpgradeTemplateAction)).
			Get("/groups", new(GroupsAction)).
			Get("/group", new(GroupAction)).
			Get("/log", new(LogAction)).
			GetPost("/update", new(UpdateAction)).
			GetPost("/test", new(TestAction)).
			GetPost("/export", new(ExportAction)).
			Get("/exportDownload", new(ExportDownloadAction)).
			GetPost("/import", new(ImportAction)).
			Post("/updateGroupOn", new(UpdateGroupOnAction)).
			Post("/deleteGroup", new(DeleteGroupAction)).
			GetPost("/createGroupPopup", new(CreateGroupPopupAction)).
			Post("/sortGroups", new(SortGroupsAction)).
			GetPost("/updateGroupPopup", new(UpdateGroupPopupAction)).
			GetPost("/createSetPopup", new(CreateSetPopupAction)).
			GetPost("/createRulePopup", new(CreateRulePopupAction)).
			Post("/sortSets", new(SortSetsAction)).
			Post("/updateSetOn", new(UpdateSetOnAction)).
			Post("/deleteSet", new(DeleteSetAction)).
			GetPost("/updateSetPopup", new(UpdateSetPopupAction)).
			Get("/setCodePopup", new(SetCodePopupAction)).
			Post("/count", new(CountAction)).
			Get("/selectPopup", new(SelectPopupAction)).
			Post("/testRegexp", new(TestRegexpAction)).

			// IP管理
			GetPost("/ipadmin", new(ipadmin.IndexAction)).
			GetPost("/ipadmin/provinces", new(ipadmin.ProvincesAction)).
			Get("/ipadmin/lists", new(ipadmin.ListsAction)).
			GetPost("/ipadmin/updateIPPopup", new(ipadmin.UpdateIPPopupAction)).
			Post("/ipadmin/deleteIP", new(ipadmin.DeleteIPAction)).
			GetPost("/ipadmin/test", new(ipadmin.TestAction)).
			EndAll()
	})
}
