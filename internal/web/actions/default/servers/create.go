package servers

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
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

	// 所有集群
	resp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClusters(this.AdminContext(), &pb.FindAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range resp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	// 服务类型
	this.Data["serverTypes"] = serverconfigs.AllServerTypes()

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name        string
	Description string
	ClusterId   int64
	GroupIds    []int64

	ServerType  string
	Addresses   string
	ServerNames string
	CertIdsJSON []byte
	Origins     string

	AccessLogIsOn bool
	WebsocketIsOn bool

	WebRoot string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入服务名称")

	if params.ClusterId <= 0 {
		this.Fail("请选择部署的集群")
	}

	// TODO 验证集群ID

	// 端口地址
	var httpConfig *serverconfigs.HTTPProtocolConfig = nil
	var httpsConfig *serverconfigs.HTTPSProtocolConfig = nil
	var tcpConfig *serverconfigs.TCPProtocolConfig = nil
	var tlsConfig *serverconfigs.TLSProtocolConfig = nil
	var unixConfig *serverconfigs.UnixProtocolConfig = nil
	var udpConfig *serverconfigs.UDPProtocolConfig = nil
	var webId int64 = 0

	switch params.ServerType {
	case serverconfigs.ServerTypeHTTPProxy, serverconfigs.ServerTypeHTTPWeb:
		listen := []*serverconfigs.NetworkAddressConfig{}
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

		listen := []*serverconfigs.NetworkAddressConfig{}
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
	case serverconfigs.ServerTypeUDPProxy:
		// 在DEMO模式下不能创建
		if teaconst.IsDemoMode {
			this.Fail("DEMO模式下不能创建UDP反向代理")
		}

		listen := []*serverconfigs.NetworkAddressConfig{}
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
	default:
		this.Fail("请选择正确的服务类型")
	}

	// 证书
	if len(params.CertIdsJSON) > 0 {
		certIds := []int64{}
		err := json.Unmarshal(params.CertIdsJSON, &certIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(certIds) > 0 {
			certRefs := []*sslconfigs.SSLCertRef{}
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
				Http2Enabled:      false,     // 默认值
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
			sslPolicyId := sslPolicyIdResp.SslPolicyId
			httpsConfig.SSLPolicyRef = &sslconfigs.SSLPolicyRef{
				IsOn:        true,
				SSLPolicyId: sslPolicyId,
			}
		}
	}

	// 域名
	if len(params.ServerNames) > 0 {
		serverNames := []*serverconfigs.ServerNameConfig{}
		err := json.Unmarshal([]byte(params.ServerNames), &serverNames)
		if err != nil {
			this.Fail("域名解析失败：" + err.Error())
		}

		// 检查域名是否已经存在
		allServerNames := serverconfigs.PlainServerNames(serverNames)
		if len(allServerNames) > 0 {
			dupResp, err := this.RPC().ServerRPC().CheckServerNameDuplicationInNodeCluster(this.AdminContext(), &pb.CheckServerNameDuplicationInNodeClusterRequest{
				ServerNames:   allServerNames,
				NodeClusterId: params.ClusterId,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if len(dupResp.DuplicatedServerNames) > 0 {
				this.Fail("域名 " + strings.Join(dupResp.DuplicatedServerNames, ", ") + " 已经被其他服务所占用，不能重复使用")
			}
		}
	}

	// 源站地址
	reverseProxyRefJSON := []byte{}
	switch params.ServerType {
	case serverconfigs.ServerTypeHTTPProxy, serverconfigs.ServerTypeTCPProxy, serverconfigs.ServerTypeUDPProxy:
		originConfigs := []*serverconfigs.OriginConfig{}
		err := json.Unmarshal([]byte(params.Origins), &originConfigs)
		if err != nil {
			this.Fail("源站地址解析失败：" + err.Error())
		}

		originRefs := []*serverconfigs.OriginRef{}
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
		reverseProxyRef := &serverconfigs.ReverseProxyRef{
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
			rootConfig := &serverconfigs.HTTPRootConfig{}
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
		webId = webResp.WebId
	}

	// 包含条件
	includeNodes := []maps.Map{}
	includeNodesJSON, err := json.Marshal(includeNodes)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 排除条件
	excludeNodes := []maps.Map{}
	excludeNodesJSON, err := json.Marshal(excludeNodes)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	req := &pb.CreateServerRequest{
		UserId:           0,
		AdminId:          this.AdminId(),
		Type:             params.ServerType,
		Name:             params.Name,
		ServerNamesJON:   []byte(params.ServerNames),
		Description:      params.Description,
		NodeClusterId:    params.ClusterId,
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
	if unixConfig != nil {
		data, err := json.Marshal(unixConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		req.UnixJSON = data
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
		webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), serverId)
		if err != nil {
			logs.Error(err)
		} else {
			// 访问日志
			if params.AccessLogIsOn {
				_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebAccessLog(this.AdminContext(), &pb.UpdateHTTPWebAccessLogRequest{
					WebId: webConfig.Id,
					AccessLogJSON: []byte(`{
			"isPrior": false,
			"isOn": true,
			"fields": [],
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
					logs.Error(err)
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
					logs.Error(err)
				} else {
					websocketId := createWebSocketResp.WebsocketId
					_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebWebsocket(this.AdminContext(), &pb.UpdateHTTPWebWebsocketRequest{
						WebId: webConfig.Id,
						WebsocketJSON: []byte(` {
				"isPrior": false,
				"isOn": true,
				"websocketId": ` + types.String(websocketId) + `
			}`),
					})
					if err != nil {
						logs.Error(err)
					}
				}
			}
		}
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "创建代理服务 %d", createResp.ServerId)

	this.Success()
}
