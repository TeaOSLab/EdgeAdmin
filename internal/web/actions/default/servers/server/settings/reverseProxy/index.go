package reverseProxy

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

// IndexAction 源站列表
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	serverTypeResp, err := this.RPC().ServerRPC().FindEnabledServerType(this.AdminContext(), &pb.FindEnabledServerTypeRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	serverType := serverTypeResp.Type

	// 当前是否有分组设置
	groupResp, err := this.RPC().ServerGroupRPC().FindEnabledServerGroupConfigInfo(this.AdminContext(), &pb.FindEnabledServerGroupConfigInfoRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasGroupConfig"] = false
	this.Data["groupSettingURL"] = ""
	switch serverType {
	case serverconfigs.ServerTypeHTTPWeb, serverconfigs.ServerTypeHTTPProxy:
		this.Data["hasGroupConfig"] = groupResp.HasHTTPReverseProxy
		if groupResp.ServerGroupId > 0 {
			this.Data["groupSettingURL"] = "/servers/groups/group/settings/httpReverseProxy?groupId=" + types.String(groupResp.ServerGroupId)
		}
	case serverconfigs.ServerTypeTCPProxy:
		this.Data["hasGroupConfig"] = groupResp.HasTCPReverseProxy
		if groupResp.ServerGroupId > 0 {
			this.Data["groupSettingURL"] = "/servers/groups/group/settings/tcpReverseProxy?groupId=" + types.String(groupResp.ServerGroupId)
		}
	case serverconfigs.ServerTypeUDPProxy:
		this.Data["hasGroupConfig"] = groupResp.HasUDPReverseProxy
		if groupResp.ServerGroupId > 0 {
			this.Data["groupSettingURL"] = "/servers/groups/group/settings/udpReverseProxy?groupId=" + types.String(groupResp.ServerGroupId)
		}
	}

	// 当前服务的配置
	reverseProxyResp, err := this.RPC().ServerRPC().FindAndInitServerReverseProxyConfig(this.AdminContext(), &pb.FindAndInitServerReverseProxyConfigRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	reverseProxyRef := &serverconfigs.ReverseProxyRef{}
	err = json.Unmarshal(reverseProxyResp.ReverseProxyRefJSON, reverseProxyRef)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["reverseProxyRef"] = reverseProxyRef

	reverseProxy := &serverconfigs.ReverseProxyConfig{}
	err = json.Unmarshal(reverseProxyResp.ReverseProxyJSON, reverseProxy)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["reverseProxyConfig"] = reverseProxy

	this.Data["serverType"] = serverType

	primaryOriginMaps := []maps.Map{}
	backupOriginMaps := []maps.Map{}
	for _, originConfig := range reverseProxy.PrimaryOrigins {
		if len(originConfig.Domains) == 0 {
			originConfig.Domains = []string{}
		}
		m := maps.Map{
			"id":      originConfig.Id,
			"weight":  originConfig.Weight,
			"addr":    originConfig.Addr.Protocol.String() + "://" + originConfig.Addr.Host + ":" + originConfig.Addr.PortRange,
			"name":    originConfig.Name,
			"isOn":    originConfig.IsOn,
			"domains": originConfig.Domains,
		}
		primaryOriginMaps = append(primaryOriginMaps, m)
	}
	for _, originConfig := range reverseProxy.BackupOrigins {
		if len(originConfig.Domains) == 0 {
			originConfig.Domains = []string{}
		}
		m := maps.Map{
			"id":      originConfig.Id,
			"weight":  originConfig.Weight,
			"addr":    originConfig.Addr.Protocol.String() + "://" + originConfig.Addr.Host + ":" + originConfig.Addr.PortRange,
			"name":    originConfig.Name,
			"isOn":    originConfig.IsOn,
			"domains": originConfig.Domains,
		}
		backupOriginMaps = append(backupOriginMaps, m)
	}
	this.Data["primaryOrigins"] = primaryOriginMaps
	this.Data["backupOrigins"] = backupOriginMaps

	this.Show()
}
