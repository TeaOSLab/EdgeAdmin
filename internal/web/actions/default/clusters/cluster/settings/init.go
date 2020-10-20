package settings

import (
	clusters "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(clusters.NewClusterHelper()).
			Prefix("/clusters/cluster/settings").
			GetPost("", new(IndexAction)).
			GetPost("/health", new(HealthAction)).
			GetPost("/healthRunPopup", new(HealthRunPopupAction)).
			EndAll()
	})
}
