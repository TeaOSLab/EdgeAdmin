package httpReverseProxy

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/groups/group/servergrouputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

// IndexAction 源站列表
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
}) {
	_, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "httpReverseProxy")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["serverType"] = "httpProxy"

	reverseProxyResp, err := this.RPC().ServerGroupRPC().FindAndInitServerGroupHTTPReverseProxyConfig(this.AdminContext(), &pb.FindAndInitServerGroupHTTPReverseProxyConfigRequest{ServerGroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var reverseProxyRef = &serverconfigs.ReverseProxyRef{}
	err = json.Unmarshal(reverseProxyResp.ReverseProxyRefJSON, reverseProxyRef)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["reverseProxyRef"] = reverseProxyRef

	var reverseProxy = serverconfigs.NewReverseProxyConfig()
	err = json.Unmarshal(reverseProxyResp.ReverseProxyJSON, reverseProxy)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["reverseProxyConfig"] = reverseProxy

	var primaryOriginMaps = []maps.Map{}
	var backupOriginMaps = []maps.Map{}
	for _, originConfig := range reverseProxy.PrimaryOrigins {
		if len(originConfig.Domains) == 0 {
			originConfig.Domains = []string{}
		}
		var m = maps.Map{
			"id":           originConfig.Id,
			"weight":       originConfig.Weight,
			"addr":         originConfig.AddrSummary(),
			"isOSS":        originConfig.IsOSS(),
			"name":         originConfig.Name,
			"isOn":         originConfig.IsOn,
			"domains":      originConfig.Domains,
			"hasCert":      originConfig.Cert != nil,
			"host":         originConfig.RequestHost,
			"followPort":   originConfig.FollowPort,
			"http2Enabled": originConfig.HTTP2Enabled,
		}
		primaryOriginMaps = append(primaryOriginMaps, m)
	}
	for _, originConfig := range reverseProxy.BackupOrigins {
		if len(originConfig.Domains) == 0 {
			originConfig.Domains = []string{}
		}
		var m = maps.Map{
			"id":           originConfig.Id,
			"weight":       originConfig.Weight,
			"addr":         originConfig.AddrSummary(),
			"isOSS":        originConfig.IsOSS(),
			"name":         originConfig.Name,
			"isOn":         originConfig.IsOn,
			"domains":      originConfig.Domains,
			"hasCert":      originConfig.Cert != nil,
			"host":         originConfig.RequestHost,
			"followPort":   originConfig.FollowPort,
			"http2Enabled": originConfig.HTTP2Enabled,
		}
		backupOriginMaps = append(backupOriginMaps, m)
	}
	this.Data["primaryOrigins"] = primaryOriginMaps
	this.Data["backupOrigins"] = backupOriginMaps

	this.Show()
}
