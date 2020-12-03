package settings

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/dns"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/toa"
	clusters "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeNode)).
			Helper(clusters.NewClusterHelper()).
			Prefix("/clusters/cluster/settings").
			GetPost("", new(IndexAction)).
			GetPost("/health", new(HealthAction)).
			GetPost("/healthRunPopup", new(HealthRunPopupAction)).

			// DNS
			Prefix("/clusters/cluster/settings/dns").
			GetPost("", new(dns.IndexAction)).

			// TOA
			Prefix("/clusters/cluster/settings/toa").
			GetPost("", new(toa.IndexAction)).

			EndAll()
	})
}
