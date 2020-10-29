package groups

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
			Prefix("/servers/components/groups").
			Get("", new(IndexAction)).
			GetPost("/createPopup", new(CreatePopupAction)).
			GetPost("/updatePopup", new(UpdatePopupAction)).
			GetPost("/selectPopup", new(SelectPopupAction)).
			Post("/delete", new(DeleteAction)).
			Post("/sort", new(SortAction)).
			EndAll()
	})
}
