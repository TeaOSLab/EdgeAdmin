package reverseProxy

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
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
	ServerId int64
}) {
	server, _, isOk := serverutils.FindServer(&this.ParentAction, params.ServerId)
	if !isOk {
		return
	}
	this.Data["serverType"] = server.Type
	this.Data["reverseProxyId"] = server.ReverseProxyId

	isOn := false
	if server.ReverseProxyId > 0 {
		reverseProxyResp, err := this.RPC().ReverseProxyRPC().FindEnabledReverseProxy(this.AdminContext(), &pb.FindEnabledReverseProxyRequest{ReverseProxyId: server.ReverseProxyId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		reverseProxy := reverseProxyResp.ReverseProxy
		if reverseProxy == nil {
			// TODO 应该在界面上提示用户开启
			this.ErrorPage(errors.New("reverse proxy should not be nil"))
			return
		}
		isOn = true

		primaryOrigins := []*serverconfigs.OriginServerConfig{}
		backupOrigins := []*serverconfigs.OriginServerConfig{}
		if len(reverseProxy.PrimaryOriginsJSON) > 0 {
			err = json.Unmarshal(reverseProxy.PrimaryOriginsJSON, &primaryOrigins)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
		if len(reverseProxy.BackupOriginsJSON) > 0 {
			err = json.Unmarshal(reverseProxy.BackupOriginsJSON, &backupOrigins)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		primaryOriginMaps := []maps.Map{}
		backupOriginMaps := []maps.Map{}
		for _, originConfig := range primaryOrigins {
			m := maps.Map{
				"id":     originConfig.Id,
				"weight": originConfig.Weight,
				"addr":   originConfig.Addr.Protocol.String() + "://" + originConfig.Addr.Host + ":" + originConfig.Addr.PortRange,
			}
			primaryOriginMaps = append(primaryOriginMaps, m)
		}
		for _, originConfig := range backupOrigins {
			m := maps.Map{
				"id":     originConfig.Id,
				"weight": originConfig.Weight,
				"addr":   originConfig.Addr.Protocol.String() + "://" + originConfig.Addr.Host + ":" + originConfig.Addr.PortRange,
			}
			backupOriginMaps = append(backupOriginMaps, m)
		}
		this.Data["primaryOrigins"] = primaryOriginMaps
		this.Data["backupOrigins"] = backupOriginMaps
	}
	this.Data["isOn"] = isOn

	this.Show()
}
