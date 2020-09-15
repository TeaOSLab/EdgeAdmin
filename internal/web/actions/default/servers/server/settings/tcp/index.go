package tcp

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

// TCP设置
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("tcp")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	server, config, isOk := serverutils.FindServer(&this.ParentAction, params.ServerId)
	if !isOk {
		return
	}
	if config.TCP == nil {
		this.ErrorPage(errors.New("there is no tcp setting"))
		return
	}

	if config.TCP.Listen == nil {
		config.TCP.Listen = []*serverconfigs.NetworkAddressConfig{}
	}

	this.Data["serverType"] = server.Type
	this.Data["addresses"] = config.TCP.Listen

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId   int64
	ServerType string
	Addresses  string

	Must *actions.Must
}) {
	serverId := params.ServerId

	server, _, isOk := serverutils.FindServer(&this.ParentAction, params.ServerId)
	if !isOk {
		return
	}

	addresses := []*serverconfigs.NetworkAddressConfig{}
	err := json.Unmarshal([]byte(params.Addresses), &addresses)
	if err != nil {
		this.Fail("端口地址解析失败：" + err.Error())
	}

	switch server.Type {
	case serverconfigs.ServerTypeHTTPProxy, serverconfigs.ServerTypeHTTPWeb:
		var httpConfig = &serverconfigs.HTTPProtocolConfig{}
		if len(server.HttpJSON) > 0 {
			err = json.Unmarshal(server.HttpJSON, httpConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			httpConfig.Listen = []*serverconfigs.NetworkAddressConfig{}
		} else {
			httpConfig.IsOn = true
		}

		var httpsConfig = &serverconfigs.HTTPSProtocolConfig{}
		if len(server.HttpsJSON) > 0 {
			err = json.Unmarshal(server.HttpsJSON, httpsConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			httpsConfig.Listen = []*serverconfigs.NetworkAddressConfig{}
		} else {
			httpsConfig.IsOn = true
		}

		for _, addr := range addresses {
			switch addr.Protocol.Primary() {
			case serverconfigs.ProtocolHTTP:
				httpConfig.AddListen(addr)
			case serverconfigs.ProtocolHTTPS:
				httpsConfig.AddListen(addr)
			}
		}

		httpData, err := json.Marshal(httpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		httpsData, err := json.Marshal(httpsConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().ServerRPC().UpdateServerHTTP(this.AdminContext(), &pb.UpdateServerHTTPRequest{
			ServerId: serverId,
			Config:   httpData,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().ServerRPC().UpdateServerHTTPS(this.AdminContext(), &pb.UpdateServerHTTPSRequest{
			ServerId: serverId,
			Config:   httpsData,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	case serverconfigs.ServerTypeTCPProxy:
		tcpProxy := &serverconfigs.TCPProtocolConfig{}
		if len(server.TcpJSON) > 0 {
			err = json.Unmarshal(server.TcpJSON, tcpProxy)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			tcpProxy.Listen = []*serverconfigs.NetworkAddressConfig{}
		} else {
			tcpProxy.IsOn = true
		}

		tlsProxy := &serverconfigs.TLSProtocolConfig{}
		if len(server.TlsJSON) > 0 {
			err = json.Unmarshal(server.TlsJSON, tlsProxy)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			tlsProxy.Listen = []*serverconfigs.NetworkAddressConfig{}
		} else {
			tlsProxy.IsOn = true
		}

		for _, addr := range addresses {
			switch addr.Protocol.Primary() {
			case serverconfigs.ProtocolTCP:
				tcpProxy.AddListen(addr)
			case serverconfigs.ProtocolTLS:
				tlsProxy.AddListen(addr)
			}
		}

		tcpData, err := json.Marshal(tcpProxy)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		tlsData, err := json.Marshal(tlsProxy)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		_, err = this.RPC().ServerRPC().UpdateServerTCP(this.AdminContext(), &pb.UpdateServerTCPRequest{
			ServerId: serverId,
			Config:   tcpData,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		_, err = this.RPC().ServerRPC().UpdateServerTLS(this.AdminContext(), &pb.UpdateServerTLSRequest{
			ServerId: serverId,
			Config:   tlsData,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	case serverconfigs.ServerTypeUnixProxy:
		unixConfig := &serverconfigs.UnixProtocolConfig{}
		if len(server.UnixJSON) > 0 {
			err = json.Unmarshal(server.UnixJSON, unixConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			unixConfig.Listen = []*serverconfigs.NetworkAddressConfig{}
		}
		for _, addr := range addresses {
			switch addr.Protocol.Primary() {
			case serverconfigs.ProtocolUnix:
				unixConfig.AddListen(addr)
			}
		}
		unixData, err := json.Marshal(unixConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().ServerRPC().UpdateServerUnix(this.AdminContext(), &pb.UpdateServerUnixRequest{
			ServerId: serverId,
			Config:   unixData,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	case serverconfigs.ServerTypeUDPProxy:
		udpConfig := &serverconfigs.UDPProtocolConfig{}
		if len(server.UdpJSON) > 0 {
			err = json.Unmarshal(server.UdpJSON, udpConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			udpConfig.Listen = []*serverconfigs.NetworkAddressConfig{}
		}
		for _, addr := range addresses {
			switch addr.Protocol.Primary() {
			case serverconfigs.ProtocolUDP:
				udpConfig.AddListen(addr)
			}
		}
		udpData, err := json.Marshal(udpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().ServerRPC().UpdateServerUDP(this.AdminContext(), &pb.UpdateServerUDPRequest{
			ServerId: serverId,
			Config:   udpData,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Success()
}
