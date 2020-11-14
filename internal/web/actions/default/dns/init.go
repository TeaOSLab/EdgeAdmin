package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/clusters"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/issues"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/providers"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Prefix("/dns").
			Get("", new(IndexAction)).
			GetPost("/updateClusterPopup", new(UpdateClusterPopupAction)).
			Post("/providerOptions", new(ProviderOptionsAction)).
			Post("/domainOptions", new(DomainOptionsAction)).

			// 集群
			Prefix("/dns/clusters").
			Get("/cluster", new(clusters.ClusterAction)).
			Post("/sync", new(clusters.SyncAction)).

			// 服务商
			Prefix("/dns/providers").
			Data("teaSubMenu", "provider").
			Get("", new(providers.IndexAction)).
			GetPost("/createPopup", new(providers.CreatePopupAction)).
			GetPost("/updatePopup", new(providers.UpdatePopupAction)).
			Post("/delete", new(providers.DeleteAction)).
			Get("/provider", new(providers.ProviderAction)).
			EndData().

			// 域名
			Prefix("/dns/domains").
			Data("teaSubMenu", "provider").
			GetPost("/createPopup", new(domains.CreatePopupAction)).
			GetPost("/updatePopup", new(domains.UpdatePopupAction)).
			Post("/delete", new(domains.DeleteAction)).
			Post("/sync", new(domains.SyncAction)).
			Get("/routesPopup", new(domains.RoutesPopupAction)).
			EndData().

			// 问题修复
			Prefix("/dns/issues").
			Data("teaSubMenu", "issue").
			GetPost("", new(issues.IndexAction)).
			GetPost("/updateNodePopup", new(issues.UpdateNodePopupAction)).
			EndData().

			EndAll()
	})
}
