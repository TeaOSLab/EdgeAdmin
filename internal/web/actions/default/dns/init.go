package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
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
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeDNS)).
			Helper(new(Helper)).
			Prefix("/dns").
			Data("teaSubMenu", "cluster").
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
			Post("/syncDomains", new(providers.SyncDomainsAction)).
			EndData().

			// 域名
			Prefix("/dns/domains").
			Data("teaSubMenu", "provider").
			GetPost("/createPopup", new(domains.CreatePopupAction)).
			GetPost("/updatePopup", new(domains.UpdatePopupAction)).
			Post("/delete", new(domains.DeleteAction)).
			Post("/recover", new(domains.RecoverAction)).
			Post("/sync", new(domains.SyncAction)).
			Get("/routesPopup", new(domains.RoutesPopupAction)).
			GetPost("/selectPopup", new(domains.SelectPopupAction)).
			Get("/clustersPopup", new(domains.ClustersPopupAction)).
			Get("/nodesPopup", new(domains.NodesPopupAction)).
			Get("/serversPopup", new(domains.ServersPopupAction)).
			EndData().

			// 问题修复
			Prefix("/dns/issues").
			Data("teaSubMenu", "issue").
			GetPost("", new(issues.IndexAction)).
			GetPost("/updateNodePopup", new(issues.UpdateNodePopupAction)).
			Post("/syncDomain", new(issues.SyncDomainAction)).
			EndData().
			EndAll()
	})
}
