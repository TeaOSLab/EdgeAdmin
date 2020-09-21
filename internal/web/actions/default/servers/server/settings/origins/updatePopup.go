package origins

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"strings"
)

// 修改源站
type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	ServerId       int64
	ReverseProxyId int64
	OriginType     string
	OriginId       int64
}) {
	this.Data["originType"] = params.OriginType
	this.Data["reverseProxyId"] = params.ReverseProxyId
	this.Data["originId"] = params.OriginId

	serverTypeResp, err := this.RPC().ServerRPC().FindEnabledServerType(this.AdminContext(), &pb.FindEnabledServerTypeRequest{
		ServerId: params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["serverType"] = serverTypeResp.Type

	// 源站信息
	originResp, err := this.RPC().OriginServerRPC().FindEnabledOriginServerConfig(this.AdminContext(), &pb.FindEnabledOriginServerConfigRequest{OriginId: params.OriginId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	configData := originResp.OriginJSON
	config := &serverconfigs.OriginServerConfig{}
	err = json.Unmarshal(configData, config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["origin"] = maps.Map{
		"id":       config.Id,
		"protocol": config.Addr.Protocol,
		"addr":     config.Addr.Host + ":" + config.Addr.PortRange,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	OriginType string
	OriginId   int64

	ReverseProxyId int64
	Protocol       string
	Addr           string

	Must *actions.Must
}) {
	params.Must.
		Field("addr", params.Addr).
		Require("请输入源站地址")

	addr := regexp.MustCompile(`\s+`).ReplaceAllString(params.Addr, "")
	portIndex := strings.LastIndex(params.Addr, ":")
	if portIndex < 0 {
		this.Fail("地址中需要带有端口")
	}
	host := addr[:portIndex]
	port := addr[portIndex+1:]

	_, err := this.RPC().OriginServerRPC().UpdateOriginServer(this.AdminContext(), &pb.UpdateOriginServerRequest{
		OriginId: params.OriginId,
		Name:     "",
		Addr: &pb.NetworkAddress{
			Protocol:  params.Protocol,
			Host:      host,
			PortRange: port,
		},
		Description: "",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	originConfigResp, err := this.RPC().OriginServerRPC().FindEnabledOriginServerConfig(this.AdminContext(), &pb.FindEnabledOriginServerConfigRequest{OriginId: params.OriginId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	originConfigData := originConfigResp.OriginJSON
	var originConfig = &serverconfigs.OriginServerConfig{}
	err = json.Unmarshal(originConfigData, originConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 查找反向代理信息
	reverseProxyResp, err := this.RPC().ReverseProxyRPC().FindEnabledReverseProxy(this.AdminContext(), &pb.FindEnabledReverseProxyRequest{ReverseProxyId: params.ReverseProxyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	reverseProxy := reverseProxyResp.ReverseProxy
	if reverseProxy == nil {
		this.ErrorPage(errors.New("reverse proxy should not be nil"))
		return
	}

	origins := []*serverconfigs.OriginServerConfig{}
	switch params.OriginType {
	case "primary":
		if len(reverseProxy.PrimaryOriginsJSON) > 0 {
			err = json.Unmarshal(reverseProxy.PrimaryOriginsJSON, &origins)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	case "backup":
		if len(reverseProxy.BackupOriginsJSON) > 0 {
			err = json.Unmarshal(reverseProxy.BackupOriginsJSON, &origins)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	}

	for index, origin := range origins {
		if origin.Id == params.OriginId {
			origins[index] = originConfig
		}
	}

	// 保存
	originsData, err := json.Marshal(origins)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	switch params.OriginType {
	case "primary":
		_, err = this.RPC().ReverseProxyRPC().UpdateReverseProxyPrimaryOrigins(this.AdminContext(), &pb.UpdateReverseProxyPrimaryOriginsRequest{
			ReverseProxyId: params.ReverseProxyId,
			OriginsJSON:    originsData,
		})
	case "backup":
		_, err = this.RPC().ReverseProxyRPC().UpdateReverseProxyBackupOrigins(this.AdminContext(), &pb.UpdateReverseProxyBackupOriginsRequest{
			ReverseProxyId: params.ReverseProxyId,
			OriginsJSON:    originsData,
		})
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
