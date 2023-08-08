package origins

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	"net/url"
	"regexp"
	"strings"
)

// AddPopupAction 添加源站
type AddPopupAction struct {
	actionutils.ParentAction
}

func (this *AddPopupAction) RunGet(params struct {
	ServerId       int64
	ServerType     string
	ReverseProxyId int64
	OriginType     string
}) {
	this.Data["reverseProxyId"] = params.ReverseProxyId
	this.Data["originType"] = params.OriginType

	var serverType string
	if params.ServerId > 0 {
		serverTypeResp, err := this.RPC().ServerRPC().FindEnabledServerType(this.AdminContext(), &pb.FindEnabledServerTypeRequest{ServerId: params.ServerId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		serverType = serverTypeResp.Type
	} else {
		serverType = params.ServerType
	}
	this.Data["serverType"] = serverType

	// 是否为HTTP
	this.Data["isHTTP"] = serverType == "httpProxy" || serverType == "httpWeb"

	// OSS
	this.getOSSHook()

	this.Show()
}

func (this *AddPopupAction) RunPost(params struct {
	OriginType string

	ReverseProxyId int64
	Weight         int32
	Protocol       string
	Addr           string
	Name           string

	ConnTimeout  int
	ReadTimeout  int
	MaxConns     int32
	MaxIdleConns int32
	IdleTimeout  int

	CertIdsJSON []byte

	DomainsJSON  []byte
	Host         string
	FollowPort   bool
	Http2Enabled bool

	Description string
	IsOn        bool

	Must *actions.Must
}) {
	ossConfig, goNext, err := this.postOSSHook(params.Protocol)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if !goNext {
		return
	}

	// 初始化
	var pbAddr = &pb.NetworkAddress{
		Protocol: params.Protocol,
	}
	var connTimeoutJSON []byte
	var readTimeoutJSON []byte
	var idleTimeoutJSON []byte
	var certRefJSON []byte

	var ossJSON []byte = nil
	if ossConfig != nil { // OSS
		ossJSON, err = json.Marshal(ossConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		err = ossConfig.Init()
		if err != nil {
			this.Fail("校验OSS配置时出错：" + err.Error())
			return
		}
	} else { // 普通源站
		params.Must.
			Field("addr", params.Addr).
			Require("请输入源站地址")

		var addr = params.Addr

		// 是否是完整的地址
		if (params.Protocol == "http" || params.Protocol == "https") && regexp.MustCompile(`^(http|https)://`).MatchString(addr) {
			u, err := url.Parse(addr)
			if err == nil {
				addr = u.Host
			}
		}

		addr = strings.ReplaceAll(addr, "：", ":")
		addr = regexp.MustCompile(`\s+`).ReplaceAllString(addr, "")
		var portIndex = strings.LastIndex(addr, ":")
		if portIndex < 0 {
			if params.Protocol == "http" {
				addr += ":80"
			} else if params.Protocol == "https" {
				addr += ":443"
			} else {
				this.FailField("addr", "源站地址中需要带有端口")
			}
			portIndex = strings.LastIndex(addr, ":")
		}
		var host = addr[:portIndex]
		var port = addr[portIndex+1:]

		// 检查端口号
		if port == "0" {
			this.FailField("addr", "源站端口号不能为0")
		}
		if !configutils.HasVariables(port) {
			// 必须是整数
			if !regexp.MustCompile(`^\d+$`).MatchString(port) {
				this.FailField("addr", "源站端口号只能为整数")
			}
			var portInt = types.Int(port)
			if portInt == 0 {
				this.FailField("addr", "源站端口号不能为0")
			}
			if portInt > 65535 {
				this.FailField("addr", "源站端口号不能大于65535")
			}
		}

		connTimeoutJSON, err = (&shared.TimeDuration{
			Count: int64(params.ConnTimeout),
			Unit:  shared.TimeDurationUnitSecond,
		}).AsJSON()
		if err != nil {
			this.ErrorPage(err)
			return
		}

		readTimeoutJSON, err = (&shared.TimeDuration{
			Count: int64(params.ReadTimeout),
			Unit:  shared.TimeDurationUnitSecond,
		}).AsJSON()
		if err != nil {
			this.ErrorPage(err)
			return
		}

		idleTimeoutJSON, err = (&shared.TimeDuration{
			Count: int64(params.IdleTimeout),
			Unit:  shared.TimeDurationUnitSecond,
		}).AsJSON()
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 证书
		var certIds = []int64{}
		if len(params.CertIdsJSON) > 0 {
			err = json.Unmarshal(params.CertIdsJSON, &certIds)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		if len(certIds) > 0 {
			var certId = certIds[0]
			if certId > 0 {
				var certRef = &sslconfigs.SSLCertRef{
					IsOn:   true,
					CertId: certId,
				}
				certRefJSON, err = json.Marshal(certRef)
				if err != nil {
					this.ErrorPage(err)
					return
				}
			}
		}

		pbAddr = &pb.NetworkAddress{
			Protocol:  params.Protocol,
			Host:      host,
			PortRange: port,
		}
	}

	// 专属域名
	var domains = []string{}
	if len(params.DomainsJSON) > 0 {
		err = json.Unmarshal(params.DomainsJSON, &domains)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 去除可能误加的斜杠
		for index, domain := range domains {
			domains[index] = strings.TrimSuffix(domain, "/")
		}
	}

	createResp, err := this.RPC().OriginRPC().CreateOrigin(this.AdminContext(), &pb.CreateOriginRequest{
		Name:            params.Name,
		Addr:            pbAddr,
		OssJSON:         ossJSON,
		Description:     params.Description,
		Weight:          params.Weight,
		IsOn:            params.IsOn,
		ConnTimeoutJSON: connTimeoutJSON,
		ReadTimeoutJSON: readTimeoutJSON,
		IdleTimeoutJSON: idleTimeoutJSON,
		MaxConns:        params.MaxConns,
		MaxIdleConns:    params.MaxIdleConns,
		CertRefJSON:     certRefJSON,
		Domains:         domains,
		Host:            params.Host,
		FollowPort:      params.FollowPort,
		Http2Enabled:    params.Http2Enabled,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	originId := createResp.OriginId
	originRef := &serverconfigs.OriginRef{
		IsOn:     true,
		OriginId: originId,
	}

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

	origins := []*serverconfigs.OriginRef{}
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
	origins = append(origins, originRef)
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

	// 日志
	defer this.CreateLogInfo(codes.ServerOrigin_LogCreateOrigin, originId)

	this.Success()
}
