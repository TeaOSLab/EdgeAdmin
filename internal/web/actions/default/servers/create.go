package servers

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "server", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
	// 审核中的数量
	countAuditingResp, err := this.RPC().ServerRPC().CountAllEnabledServersMatch(this.AdminContext(), &pb.CountAllEnabledServersMatchRequest{
		AuditingFlag: 1,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countAuditing"] = countAuditingResp.Count

	// 服务类型
	this.Data["serverTypes"] = serverconfigs.AllServerTypes()

	// 检查是否有用户
	countUsersResp, err := this.RPC().UserRPC().CountAllEnabledUsers(this.AdminContext(), &pb.CountAllEnabledUsersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasUsers"] = countUsersResp.Count > 0

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name        string
	Description string

	UserId     int64
	UserPlanId int64
	ClusterId  int64

	GroupIds []int64

	ServerType  string
	Addresses   string
	ServerNames []byte
	CertIdsJSON []byte
	Origins     string

	AccessLogIsOn  bool
	WebsocketIsOn  bool
	CacheIsOn      bool
	WafIsOn        bool
	RemoteAddrIsOn bool
	StatIsOn       bool

	WebRoot string

	Must *actions.Must
}) {
	var clusterId = params.ClusterId

	// 用户
	var userId = params.UserId
	if userId > 0 {
		clusterIdResp, err := this.RPC().UserRPC().FindUserNodeClusterId(this.AdminContext(), &pb.FindUserNodeClusterIdRequest{UserId: userId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		clusterId = clusterIdResp.NodeClusterId
		if clusterId <= 0 {
			this.Fail("请选择部署的集群")
		}
	}

	// 套餐
	var userPlanId = params.UserPlanId

	// 端口地址
	var httpConfig *serverconfigs.HTTPProtocolConfig = nil
	var httpsConfig *serverconfigs.HTTPSProtocolConfig = nil
	var tcpConfig *serverconfigs.TCPProtocolConfig = nil
	var tlsConfig *serverconfigs.TLSProtocolConfig = nil
	var udpConfig *serverconfigs.UDPProtocolConfig = nil
	var webId int64 = 0

	switch params.ServerType {
	case serverconfigs.ServerTypeHTTPProxy, serverconfigs.ServerTypeHTTPWeb:
		var listen = []*serverconfigs.NetworkAddressConfig{}
		err := json.Unmarshal([]byte(params.Addresses), &listen)
		if err != nil {
			this.Fail("端口地址解析失败：" + err.Error())
		}
		if len(listen) == 0 {
			this.Fail("至少需要绑定一个端口")
		}

		for _, addr := range listen {
			switch addr.Protocol.Primary() {
			case serverconfigs.ProtocolHTTP:
				if httpConfig == nil {
					httpConfig = &serverconfigs.HTTPProtocolConfig{
						BaseProtocol: serverconfigs.BaseProtocol{
							IsOn: true,
						},
					}
				}
				httpConfig.AddListen(addr)
			case serverconfigs.ProtocolHTTPS:
				if httpsConfig == nil {
					httpsConfig = &serverconfigs.HTTPSProtocolConfig{
						BaseProtocol: serverconfigs.BaseProtocol{
							IsOn: true,
						},
					}
				}
				httpsConfig.AddListen(addr)
			}
		}
	case serverconfigs.ServerTypeTCPProxy:
		// 在DEMO模式下不能创建
		if teaconst.IsDemoMode {
			this.Fail("DEMO模式下不能创建TCP反向代理")
		}

		var listen = []*serverconfigs.NetworkAddressConfig{}
		err := json.Unmarshal([]byte(params.Addresses), &listen)
		if err != nil {
			this.Fail("端口地址解析失败：" + err.Error())
		}
		if len(listen) == 0 {
			this.Fail("至少需要绑定一个端口")
		}

		for _, addr := range listen {
			switch addr.Protocol.Primary() {
			case serverconfigs.ProtocolTCP:
				if tcpConfig == nil {
					tcpConfig = &serverconfigs.TCPProtocolConfig{
						BaseProtocol: serverconfigs.BaseProtocol{
							IsOn: true,
						},
					}
				}
				tcpConfig.AddListen(addr)
			case serverconfigs.ProtocolTLS:
				if tlsConfig == nil {
					tlsConfig = &serverconfigs.TLSProtocolConfig{
						BaseProtocol: serverconfigs.BaseProtocol{
							IsOn: true,
						},
					}
				}
				tlsConfig.AddListen(addr)
			}
		}

		if len(params.Name) == 0 {
			params.Name = "TCP负载均衡"
		}
	case serverconfigs.ServerTypeUDPProxy:
		// 在DEMO模式下不能创建
		if teaconst.IsDemoMode {
			this.Fail("DEMO模式下不能创建UDP反向代理")
		}

		var listen = []*serverconfigs.NetworkAddressConfig{}
		err := json.Unmarshal([]byte(params.Addresses), &listen)
		if err != nil {
			this.Fail("端口地址解析失败：" + err.Error())
		}
		if len(listen) == 0 {
			this.Fail("至少需要绑定一个端口")
		}

		for _, addr := range listen {
			switch addr.Protocol.Primary() {
			case serverconfigs.ProtocolUDP:
				if udpConfig == nil {
					udpConfig = &serverconfigs.UDPProtocolConfig{
						BaseProtocol: serverconfigs.BaseProtocol{
							IsOn: true,
						},
					}
				}
				udpConfig.AddListen(addr)
			}
		}

		if len(params.Name) == 0 {
			params.Name = "UDP负载均衡"
		}
	default:
		this.Fail("请选择正确的服务类型")
	}

	// 证书
	if len(params.CertIdsJSON) > 0 {
		var certIds = []int64{}
		err := json.Unmarshal(params.CertIdsJSON, &certIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(certIds) > 0 {
			var certRefs = []*sslconfigs.SSLCertRef{}
			for _, certId := range certIds {
				certRefs = append(certRefs, &sslconfigs.SSLCertRef{
					IsOn:   true,
					CertId: certId,
				})
			}
			certRefsJSON, err := json.Marshal(certRefs)
			if err != nil {
				this.ErrorPage(err)
				return
			}

			sslPolicyIdResp, err := this.RPC().SSLPolicyRPC().CreateSSLPolicy(this.AdminContext(), &pb.CreateSSLPolicyRequest{
				Http2Enabled:      false, // 默认值
				Http3Enabled:      false,
				MinVersion:        "TLS 1.1", // 默认值
				SslCertsJSON:      certRefsJSON,
				HstsJSON:          nil,
				ClientAuthType:    0,
				ClientCACertsJSON: nil,
				CipherSuites:      nil,
				CipherSuitesIsOn:  false,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			var sslPolicyId = sslPolicyIdResp.SslPolicyId
			httpsConfig.SSLPolicyRef = &sslconfigs.SSLPolicyRef{
				IsOn:        true,
				SSLPolicyId: sslPolicyId,
			}
		}
	}

	// 域名
	var serverNames = []*serverconfigs.ServerNameConfig{}
	if len(params.ServerNames) > 0 {
		err := json.Unmarshal(params.ServerNames, &serverNames)
		if err != nil {
			this.Fail("域名解析失败：" + err.Error())
		}

		// 检查域名是否已经存在
		var allServerNames = serverconfigs.PlainServerNames(serverNames)
		if len(allServerNames) > 0 {
			// 指定默认名称
			if len(params.Name) == 0 {
				params.Name = allServerNames[0]
			}

			dupResp, err := this.RPC().ServerRPC().CheckServerNameDuplicationInNodeCluster(this.AdminContext(), &pb.CheckServerNameDuplicationInNodeClusterRequest{
				ServerNames:   allServerNames,
				NodeClusterId: clusterId,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if len(dupResp.DuplicatedServerNames) > 0 {
				this.Fail("域名 " + strings.Join(dupResp.DuplicatedServerNames, ", ") + " 已经被其他网站所占用，不能重复使用")
			}
		}
	}
	if params.ServerType == serverconfigs.ServerTypeHTTPProxy && len(serverNames) == 0 {
		this.FailField("emptyDomain", "请输入添加至少一个域名")
	}

	// 源站地址
	var reverseProxyRefJSON = []byte{}
	switch params.ServerType {
	case serverconfigs.ServerTypeHTTPProxy, serverconfigs.ServerTypeTCPProxy, serverconfigs.ServerTypeUDPProxy:
		var originConfigs = []*serverconfigs.OriginConfig{}
		err := json.Unmarshal([]byte(params.Origins), &originConfigs)
		if err != nil {
			this.Fail("源站地址解析失败：" + err.Error())
		}
		if len(originConfigs) == 0 {
			this.FailField("emptyOrigin", "请添加至少一个源站地址")
		}

		var originRefs = []*serverconfigs.OriginRef{}
		for _, originConfig := range originConfigs {
			if originConfig.Id > 0 {
				originRefs = append(originRefs, &serverconfigs.OriginRef{
					IsOn:     true,
					OriginId: originConfig.Id,
				})
			}
		}
		originRefsJSON, err := json.Marshal(originRefs)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		resp, err := this.RPC().ReverseProxyRPC().CreateReverseProxy(this.AdminContext(), &pb.CreateReverseProxyRequest{
			SchedulingJSON:     nil,
			PrimaryOriginsJSON: originRefsJSON,
			BackupOriginsJSON:  nil,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var reverseProxyRef = &serverconfigs.ReverseProxyRef{
			IsOn:           true,
			ReverseProxyId: resp.ReverseProxyId,
		}
		reverseProxyRefJSON, err = json.Marshal(reverseProxyRef)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// Web地址
	switch params.ServerType {
	case serverconfigs.ServerTypeHTTPWeb:
		var rootJSON []byte
		var err error
		if len(params.WebRoot) > 0 {
			var rootConfig = serverconfigs.NewHTTPRootConfig()
			rootConfig.IsOn = true
			rootConfig.Dir = params.WebRoot
			rootConfig.Indexes = []string{"index.html", "index.htm"}
			rootJSON, err = json.Marshal(rootConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		webResp, err := this.RPC().HTTPWebRPC().CreateHTTPWeb(this.AdminContext(), &pb.CreateHTTPWebRequest{RootJSON: rootJSON})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		webId = webResp.HttpWebId
	}

	// 包含条件
	var includeNodes = []maps.Map{}
	includeNodesJSON, err := json.Marshal(includeNodes)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 排除条件
	var excludeNodes = []maps.Map{}
	excludeNodesJSON, err := json.Marshal(excludeNodes)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var req = &pb.CreateServerRequest{
		UserId:           userId,
		UserPlanId:       userPlanId,
		AdminId:          this.AdminId(),
		Type:             params.ServerType,
		Name:             params.Name,
		ServerNamesJSON:  params.ServerNames,
		Description:      params.Description,
		NodeClusterId:    clusterId,
		IncludeNodesJSON: includeNodesJSON,
		ExcludeNodesJSON: excludeNodesJSON,
		WebId:            webId,
		ReverseProxyJSON: reverseProxyRefJSON,
		ServerGroupIds:   params.GroupIds,
	}
	if httpConfig != nil {
		data, err := json.Marshal(httpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		req.HttpJSON = data
	}
	if httpsConfig != nil {
		data, err := json.Marshal(httpsConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		req.HttpsJSON = data
	}
	if tcpConfig != nil {
		data, err := json.Marshal(tcpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		req.TcpJSON = data
	}
	if tlsConfig != nil {
		data, err := json.Marshal(tlsConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		req.TlsJSON = data
	}
	if udpConfig != nil {
		data, err := json.Marshal(udpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		req.UdpJSON = data
	}
	createResp, err := this.RPC().ServerRPC().CreateServer(this.AdminContext(), req)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var serverId = createResp.ServerId

	// 开启访问日志和Websocket
	if params.ServerType == serverconfigs.ServerTypeHTTPProxy {
		webConfig, findErr := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), serverId)
		if findErr != nil {
			this.ErrorPage(findErr)
			return
		}
		// 访问日志
		if params.AccessLogIsOn {
			_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebAccessLog(this.AdminContext(), &pb.UpdateHTTPWebAccessLogRequest{
				HttpWebId: webConfig.Id,
				AccessLogJSON: []byte(`{
			"isPrior": false,
			"isOn": true,
			"fields": [1, 2, 6, 7],
			"status1": true,
			"status2": true,
			"status3": true,
			"status4": true,
			"status5": true,

			"storageOnly": false,
			"storagePolicies": [],

            "firewallOnly": false
		}`),
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		// websocket
		if params.WebsocketIsOn {
			createWebSocketResp, err := this.RPC().HTTPWebsocketRPC().CreateHTTPWebsocket(this.AdminContext(), &pb.CreateHTTPWebsocketRequest{
				HandshakeTimeoutJSON: []byte(`{
					"count": 30,
					"unit": "second"
				}`),
				AllowAllOrigins:   true,
				AllowedOrigins:    nil,
				RequestSameOrigin: true,
				RequestOrigin:     "",
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			websocketId := createWebSocketResp.WebsocketId
			_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebWebsocket(this.AdminContext(), &pb.UpdateHTTPWebWebsocketRequest{
				HttpWebId: webConfig.Id,
				WebsocketJSON: []byte(`{
				"isPrior": false,
				"isOn": true,
				"websocketId": ` + types.String(websocketId) + `
			}`),
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		// cache
		if params.CacheIsOn {
			var cacheConfig = &serverconfigs.HTTPCacheConfig{
				IsPrior:         false,
				IsOn:            true,
				AddStatusHeader: true,
				PurgeIsOn:       false,
				PurgeKey:        "",
				CacheRefs:       []*serverconfigs.HTTPCacheRef{},
			}
			cacheConfigJSON, err := json.Marshal(cacheConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebCache(this.AdminContext(), &pb.UpdateHTTPWebCacheRequest{
				HttpWebId: webConfig.Id,
				CacheJSON: cacheConfigJSON,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		// waf
		if params.WafIsOn {
			var firewallRef = &firewallconfigs.HTTPFirewallRef{
				IsPrior:          false,
				IsOn:             true,
				FirewallPolicyId: 0,
			}
			firewallRefJSON, err := json.Marshal(firewallRef)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebFirewall(this.AdminContext(), &pb.UpdateHTTPWebFirewallRequest{
				HttpWebId:    webConfig.Id,
				FirewallJSON: firewallRefJSON,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		// remoteAddr
		{
			var remoteAddrConfig = &serverconfigs.HTTPRemoteAddrConfig{
				IsOn:  true,
				Value: "${rawRemoteAddr}",
				Type:  serverconfigs.HTTPRemoteAddrTypeDefault,
			}
			if params.RemoteAddrIsOn {
				remoteAddrConfig.Value = "${remoteAddr}"
				remoteAddrConfig.Type = serverconfigs.HTTPRemoteAddrTypeProxy
			}
			remoteAddrConfigJSON, err := json.Marshal(remoteAddrConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebRemoteAddr(this.AdminContext(), &pb.UpdateHTTPWebRemoteAddrRequest{
				HttpWebId:      webConfig.Id,
				RemoteAddrJSON: remoteAddrConfigJSON,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		// 统计
		if params.StatIsOn {
			var statConfig = &serverconfigs.HTTPStatRef{
				IsPrior: false,
				IsOn:    true,
			}
			statJSON, err := json.Marshal(statConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebStat(this.AdminContext(), &pb.UpdateHTTPWebStatRequest{
				HttpWebId: webConfig.Id,
				StatJSON:  statJSON,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	}

	// 创建日志
	defer this.CreateLogInfo(codes.Server_LogCreateServer, createResp.ServerId)

	this.Success()
}
