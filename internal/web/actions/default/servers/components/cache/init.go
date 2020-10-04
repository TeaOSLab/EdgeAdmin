package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/componentutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(NewHelper()).
			Helper(componentutils.NewComponentHelper()).
			Prefix("/servers/components/cache").
			Get("", new(IndexAction)).
			GetPost("/createPopup", new(CreatePopupAction)).
			Get("/policy", new(PolicyAction)).
			GetPost("/update", new(UpdateAction)).
			GetPost("/clean", new(CleanAction)).
			GetPost("/preheat", new(PreheatAction)).
			GetPost("/purge", new(PurgeAction)).
			GetPost("/stat", new(StatAction)).
			GetPost("/test", new(TestAction)).
			Post("/delete", new(DeleteAction)).
			Post("/testRead", new(TestReadAction)).
			Post("/testWrite", new(TestWriteAction)).
			EndAll()
	})
}
