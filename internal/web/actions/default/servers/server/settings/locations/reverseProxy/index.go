package reverseProxy

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

// 源站列表
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ServerId   int64
	LocationId int64
}) {
	serverTypeResp, err := this.RPC().ServerRPC().FindEnabledServerType(this.AdminContext(), &pb.FindEnabledServerTypeRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var serverType = serverTypeResp.Type

	reverseProxyResp, err := this.RPC().HTTPLocationRPC().FindAndInitHTTPLocationReverseProxyConfig(this.AdminContext(), &pb.FindAndInitHTTPLocationReverseProxyConfigRequest{LocationId: params.LocationId})
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

	this.Data["serverType"] = serverType

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
