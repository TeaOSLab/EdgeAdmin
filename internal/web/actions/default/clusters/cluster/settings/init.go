package settings

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/cache"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/dns"
	firewallActions "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/firewall-actions"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/services"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/toa"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/settings/waf"
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

			// 健康检查
			GetPost("/health", new(HealthAction)).
			GetPost("/healthRunPopup", new(HealthRunPopupAction)).

			// 缓存
			GetPost("/cache", new(cache.IndexAction)).

			// WAF
			GetPost("/waf", new(waf.IndexAction)).

			// DNS
			Prefix("/clusters/cluster/settings/dns").
			GetPost("", new(dns.IndexAction)).

			// TOA
			Prefix("/clusters/cluster/settings/toa").
			GetPost("", new(toa.IndexAction)).

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

			EndAll()
	})
}
