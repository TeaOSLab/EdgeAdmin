package reverseProxy

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
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
	serverTypeResp, err := this.RPC().ServerRPC().FindEnabledServerType(this.AdminContext(), &pb.FindEnabledServerTypeRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	serverType := serverTypeResp.Type

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

	reverseProxy := &serverconfigs.ReverseProxyConfig{}
	err = json.Unmarshal(reverseProxyResp.ReverseProxyJSON, reverseProxy)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["serverType"] = serverType
	this.Data["reverseProxyRef"] = reverseProxyRef
	this.Data["reverseProxyConfig"] = reverseProxy

	primaryOriginMaps := []maps.Map{}
	backupOriginMaps := []maps.Map{}
	for _, originConfig := range reverseProxy.PrimaryOrigins {
		m := maps.Map{
			"id":     originConfig.Id,
			"weight": originConfig.Weight,
			"addr":   originConfig.Addr.Protocol.String() + "://" + originConfig.Addr.Host + ":" + originConfig.Addr.PortRange,
		}
		primaryOriginMaps = append(primaryOriginMaps, m)
	}
	for _, originConfig := range reverseProxy.BackupOrigins {
		m := maps.Map{
			"id":     originConfig.Id,
			"weight": originConfig.Weight,
			"addr":   originConfig.Addr.Protocol.String() + "://" + originConfig.Addr.Host + ":" + originConfig.Addr.PortRange,
		}
		backupOriginMaps = append(backupOriginMaps, m)
	}
	this.Data["primaryOrigins"] = primaryOriginMaps
	this.Data["backupOrigins"] = backupOriginMaps

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId            int64
	ReverseProxyRefJSON []byte

	Must *actions.Must
}) {
	// TODO 校验配置

	_, err := this.RPC().ServerRPC().UpdateServerReverseProxy(this.AdminContext(), &pb.UpdateServerReverseProxyRequest{
		ServerId:         params.ServerId,
		ReverseProxyJSON: params.ReverseProxyRefJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
