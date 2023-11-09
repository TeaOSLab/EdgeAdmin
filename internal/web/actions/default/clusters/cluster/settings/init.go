package settings

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/cache"
	ddosProtection "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/ddos-protection"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/dns"
	firewallActions "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/firewall-actions"
	globalServerConfig "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/global-server-config"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/health"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/metrics"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/services"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/waf"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/webp"
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
			Data("teaSubMenu", "cluster").
			GetPost("", new(IndexAction)).

			// 健康检查
			GetPost("/health", new(health.IndexAction)).
			GetPost("/health/runPopup", new(health.RunPopupAction)).
			Post("/health/checkDomain", new(health.CheckDomainAction)).

			// 缓存
			GetPost("/cache", new(cache.IndexAction)).

			// WAF
			GetPost("/waf", new(waf.IndexAction)).

			// DNS
			Prefix("/clusters/cluster/settings/dns").
			GetPost("", new(dns.IndexAction)).
			Get("/records", new(dns.RecordsAction)).
			Post("/randomName", new(dns.RandomNameAction)).

			// 系统服务设置
			Prefix("/clusters/cluster/settings/services").
			GetPost("", new(services.IndexAction)).
			GetPost("/status", new(services.StatusAction)).

			// 防火墙动作
			Prefix("/clusters/cluster/settings/firewall-actions").
			Get("", new(firewallActions.IndexAction)).
			GetPost("/createPopup", new(firewallActions.CreatePopupAction)).
			GetPost("/updatePopup", new(firewallActions.UpdatePopupAction)).
			Post("/delete", new(firewallActions.DeleteAction)).

			// 指标
			Prefix("/clusters/cluster/settings/metrics").
			Get("", new(metrics.IndexAction)).
			GetPost("/createPopup", new(metrics.CreatePopupAction)).
			Post("/delete", new(metrics.DeleteAction)).

			// WebP
			Prefix("/clusters/cluster/settings/webp").
			GetPost("", new(webp.IndexAction)).

			// DDOS Protection
			Prefix("/clusters/cluster/settings/ddos-protection").
			GetPost("", new(ddosProtection.IndexAction)).
			GetPost("/status", new(ddosProtection.StatusAction)).

			// 全局服务配置
			Prefix("/clusters/cluster/settings/global-server-config").
			GetPost("", new(globalServerConfig.IndexAction)).

			//
			EndAll()
	})
}
