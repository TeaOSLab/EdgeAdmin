package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "server", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
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
	for _, cluster := range resp.Clusters {
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

	ServerType  string
	Addresses   string
	ServerNames string
	Origins     string

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

	// 配置
	serverConfig := &serverconfigs.ServerConfig{}
	serverConfig.IsOn = true
	serverConfig.Name = params.Name
	serverConfig.Description = params.Description

	// 端口地址
	switch params.ServerType {
	case serverconfigs.ServerTypeHTTPProxy, serverconfigs.ServerTypeHTTPWeb:
		listen := []*serverconfigs.NetworkAddressConfig{}
		err := json.Unmarshal([]byte(params.Addresses), &listen)
		if err != nil {
			this.Fail("端口地址解析失败：" + err.Error())
		}

		for _, addr := range listen {
			switch addr.Protocol {
			case serverconfigs.ProtocolHTTP, serverconfigs.ProtocolHTTP4, serverconfigs.ProtocolHTTP6:
				if serverConfig.HTTP == nil {
					serverConfig.HTTP = &serverconfigs.HTTPProtocolConfig{
						BaseProtocol: serverconfigs.BaseProtocol{
							IsOn: true,
						},
					}
				}
				serverConfig.HTTP.AddListen(addr)
			case serverconfigs.ProtocolHTTPS, serverconfigs.ProtocolHTTPS4, serverconfigs.ProtocolHTTPS6:
				if serverConfig.HTTPS == nil {
					serverConfig.HTTPS = &serverconfigs.HTTPSProtocolConfig{
						BaseProtocol: serverconfigs.BaseProtocol{
							IsOn: true,
						},
					}
				}
				serverConfig.HTTPS.AddListen(addr)
			}
		}
	case serverconfigs.ServerTypeTCPProxy:
		listen := []*serverconfigs.NetworkAddressConfig{}
		err := json.Unmarshal([]byte(params.Addresses), &listen)
		if err != nil {
			this.Fail("端口地址解析失败：" + err.Error())
		}

		for _, addr := range listen {
			switch addr.Protocol {
			case serverconfigs.ProtocolTCP, serverconfigs.ProtocolTCP4, serverconfigs.ProtocolTCP6:
				if serverConfig.TCP == nil {
					serverConfig.TCP = &serverconfigs.TCPProtocolConfig{
						BaseProtocol: serverconfigs.BaseProtocol{
							IsOn: true,
						},
					}
				}
				serverConfig.TCP.AddListen(addr)
			case serverconfigs.ProtocolTLS, serverconfigs.ProtocolTLS4, serverconfigs.ProtocolTLS6:
				if serverConfig.TLS == nil {
					serverConfig.TLS = &serverconfigs.TLSProtocolConfig{
						BaseProtocol: serverconfigs.BaseProtocol{
							IsOn: true,
						},
					}
				}
				serverConfig.TLS.AddListen(addr)
			}
		}
	default:
		this.Fail("请选择正确的服务类型")
	}

	// TODO 证书

	// 域名
	serverNames := []*serverconfigs.ServerNameConfig{}
	err := json.Unmarshal([]byte(params.ServerNames), &serverNames)
	if err != nil {
		this.Fail("域名解析失败：" + err.Error())
	}
	serverConfig.ServerNames = serverNames

	// 源站地址
	switch params.ServerType {
	case serverconfigs.ServerTypeHTTPProxy, serverconfigs.ServerTypeTCPProxy:
		origins := []*serverconfigs.OriginServerConfig{}
		err = json.Unmarshal([]byte(params.Origins), &origins)
		if err != nil {
			this.Fail("源站地址解析失败：" + err.Error())
		}
		serverConfig.ReverseProxy = &serverconfigs.ReverseProxyConfig{
			IsOn:    true,
			Origins: origins,
		}
	}

	// Web地址
	switch params.ServerType {
	case serverconfigs.ServerTypeHTTPWeb:
		serverConfig.Web = &serverconfigs.WebConfig{
			IsOn: true,
			Root: params.WebRoot,
		}
	}

	// 校验
	err = serverConfig.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
	}

	serverConfigJSON, err := serverConfig.AsJSON()
	if err != nil {
		this.ErrorPage(err)
		return
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

	_, err = this.RPC().ServerRPC().CreateServer(this.AdminContext(), &pb.CreateServerRequest{
		UserId:           0,
		AdminId:          this.AdminId(),
		Type:             params.ServerType,
		Name:             params.Name,
		Description:      params.Description,
		ClusterId:        params.ClusterId,
		Config:           serverConfigJSON,
		IncludeNodesJSON: includeNodesJSON,
		ExcludeNodesJSON: excludeNodesJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
