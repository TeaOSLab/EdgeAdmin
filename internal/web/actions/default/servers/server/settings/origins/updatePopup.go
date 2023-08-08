package origins

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
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
	ServerType     string
	ReverseProxyId int64
	OriginType     string
	OriginId       int64
}) {
	this.Data["originType"] = params.OriginType
	this.Data["reverseProxyId"] = params.ReverseProxyId
	this.Data["originId"] = params.OriginId

	var serverType string
	if params.ServerId > 0 {
		serverTypeResp, err := this.RPC().ServerRPC().FindEnabledServerType(this.AdminContext(), &pb.FindEnabledServerTypeRequest{
			ServerId: params.ServerId,
		})
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

	// 源站信息
	originResp, err := this.RPC().OriginRPC().FindEnabledOriginConfig(this.AdminContext(), &pb.FindEnabledOriginConfigRequest{OriginId: params.OriginId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var configData = originResp.OriginJSON
	var config = &serverconfigs.OriginConfig{}
	err = json.Unmarshal(configData, config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var connTimeout = 0
	var readTimeout = 0
	var idleTimeout = 0
	if config.ConnTimeout != nil {
		connTimeout = types.Int(config.ConnTimeout.Count)
	}

	if config.ReadTimeout != nil {
		readTimeout = types.Int(config.ReadTimeout.Count)
	}

	if config.IdleTimeout != nil {
		idleTimeout = types.Int(config.IdleTimeout.Count)
	}

	if len(config.Domains) == 0 {
		config.Domains = []string{}
	}

	// 重置数据
	if config.Cert != nil {
		config.Cert.CertData = nil
		config.Cert.KeyData = nil
	}

	var addr = ""
	if len(config.Addr.Host) > 0 && len(config.Addr.PortRange) > 0 {
		addr = config.Addr.Host + ":" + config.Addr.PortRange
	}
	this.Data["origin"] = maps.Map{
		"id":           config.Id,
		"protocol":     config.Addr.Protocol,
		"addr":         addr,
		"weight":       config.Weight,
		"name":         config.Name,
		"description":  config.Description,
		"isOn":         config.IsOn,
		"connTimeout":  connTimeout,
		"readTimeout":  readTimeout,
		"idleTimeout":  idleTimeout,
		"maxConns":     config.MaxConns,
		"maxIdleConns": config.MaxIdleConns,
		"cert":         config.Cert,
		"domains":      config.Domains,
		"host":         config.RequestHost,
		"followPort":   config.FollowPort,
		"http2Enabled": config.HTTP2Enabled,
		"oss":          config.OSS,
	}

	// OSS
	this.getOSSHook()

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

	var ossJSON []byte = nil
	var connTimeoutJSON []byte
	var readTimeoutJSON []byte
	var idleTimeoutJSON []byte
	var certRefJSON []byte
	var pbAddr = &pb.NetworkAddress{
		Protocol: params.Protocol,
	}

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
		portIndex := strings.LastIndex(addr, ":")
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
			return
		}
		if !configutils.HasVariables(port) {
			// 必须是整数
			if !regexp.MustCompile(`^\d+$`).MatchString(port) {
				this.FailField("addr", "源站端口号只能为整数")
				return
			}
			var portInt = types.Int(port)
			if portInt == 0 {
				this.FailField("addr", "源站端口号不能为0")
				return
			}
			if portInt > 65535 {
				this.FailField("addr", "源站端口号不能大于65535")
				return
			}
		}

		pbAddr = &pb.NetworkAddress{
			Protocol:  params.Protocol,
			Host:      host,
			PortRange: port,
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

	_, err = this.RPC().OriginRPC().UpdateOrigin(this.AdminContext(), &pb.UpdateOriginRequest{
		OriginId:        params.OriginId,
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

	// 日志
	defer this.CreateLogInfo(codes.ServerOrigin_LogUpdateOrigin, params.OriginId)

	this.Success()
}
