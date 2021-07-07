package origins

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net/url"
	"regexp"
	"strings"
)

// UpdatePopupAction 修改源站
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
	serverType := serverTypeResp.Type

	// 是否为HTTP
	this.Data["isHTTP"] = serverType == "httpProxy" || serverType == "httpWeb"

	// 源站信息
	originResp, err := this.RPC().OriginRPC().FindEnabledOriginConfig(this.AdminContext(), &pb.FindEnabledOriginConfigRequest{OriginId: params.OriginId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	configData := originResp.OriginJSON
	config := &serverconfigs.OriginConfig{}
	err = json.Unmarshal(configData, config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	connTimeout := 0
	readTimeout := 0
	idleTimeout := 0
	if config.ConnTimeout != nil {
		connTimeout = types.Int(config.ConnTimeout.Count)
	}

	if config.ReadTimeout != nil {
		readTimeout = types.Int(config.ReadTimeout.Count)
	}

	if config.IdleTimeout != nil {
		idleTimeout = types.Int(config.IdleTimeout.Count)
	}

	this.Data["origin"] = maps.Map{
		"id":           config.Id,
		"protocol":     config.Addr.Protocol,
		"addr":         config.Addr.Host + ":" + config.Addr.PortRange,
		"weight":       config.Weight,
		"name":         config.Name,
		"description":  config.Description,
		"isOn":         config.IsOn,
		"connTimeout":  connTimeout,
		"readTimeout":  readTimeout,
		"idleTimeout":  idleTimeout,
		"maxConns":     config.MaxConns,
		"maxIdleConns": config.MaxIdleConns,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	OriginType string
	OriginId   int64

	ReverseProxyId int64
	Protocol       string
	Addr           string
	Weight         int32
	Name           string

	ConnTimeout  int
	ReadTimeout  int
	MaxConns     int32
	MaxIdleConns int32
	IdleTimeout  int

	Description string
	IsOn        bool

	Must *actions.Must
}) {
	params.Must.
		Field("addr", params.Addr).
		Require("请输入源站地址")

	addr := params.Addr

	// 是否是完整的地址
	if (params.Protocol == "http" || params.Protocol == "https") && regexp.MustCompile(`^(http|https)://`).MatchString(addr) {
		u, err := url.Parse(addr)
		if err == nil {
			addr = u.Host
		}
	}

	addr = strings.ReplaceAll(addr, "：", ":")
	addr = regexp.MustCompile(`\s+`).ReplaceAllString(addr, "")
	portIndex := strings.LastIndex(addr, ":")
	if portIndex < 0 {
		if params.Protocol == "http" {
			addr += ":80"
		} else if params.Protocol == "https" {
			addr += ":443"
		} else {
			this.Fail("地址中需要带有端口")
		}
		portIndex = strings.LastIndex(addr, ":")
	}
	host := addr[:portIndex]
	port := addr[portIndex+1:]
	if port == "0" {
		this.Fail("端口号不能为0")
	}

	connTimeoutJSON, err := (&shared.TimeDuration{
		Count: int64(params.ConnTimeout),
		Unit:  shared.TimeDurationUnitSecond,
	}).AsJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	readTimeoutJSON, err := (&shared.TimeDuration{
		Count: int64(params.ReadTimeout),
		Unit:  shared.TimeDurationUnitSecond,
	}).AsJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	idleTimeoutJSON, err := (&shared.TimeDuration{
		Count: int64(params.IdleTimeout),
		Unit:  shared.TimeDurationUnitSecond,
	}).AsJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().OriginRPC().UpdateOrigin(this.AdminContext(), &pb.UpdateOriginRequest{
		OriginId: params.OriginId,
		Name:     params.Name,
		Addr: &pb.NetworkAddress{
			Protocol:  params.Protocol,
			Host:      host,
			PortRange: port,
		},
		Description:     params.Description,
		Weight:          params.Weight,
		IsOn:            params.IsOn,
		ConnTimeoutJSON: connTimeoutJSON,
		ReadTimeoutJSON: readTimeoutJSON,
		IdleTimeoutJSON: idleTimeoutJSON,
		MaxConns:        params.MaxConns,
		MaxIdleConns:    params.MaxIdleConns,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "修改反向代理服务 %d 的源站 %d", params.ReverseProxyId, params.OriginId)

	this.Success()
}
