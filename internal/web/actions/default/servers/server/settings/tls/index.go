package tls

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

// TLS设置
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("tls")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	tlsConfig := &serverconfigs.TLSProtocolConfig{}
	if len(server.TlsJSON) > 0 {
		err := json.Unmarshal(server.TlsJSON, tlsConfig)
		if err != nil {
			this.ErrorPage(err)
		}
	} else {
		tlsConfig.IsOn = true
	}

	// SSL配置
	var sslPolicy *sslconfigs.SSLPolicy
	if tlsConfig.SSLPolicyRef != nil && tlsConfig.SSLPolicyRef.SSLPolicyId > 0 {
		sslPolicyConfigResp, err := this.RPC().SSLPolicyRPC().FindEnabledSSLPolicyConfig(this.AdminContext(), &pb.FindEnabledSSLPolicyConfigRequest{
			SslPolicyId: tlsConfig.SSLPolicyRef.SSLPolicyId,
			IgnoreData:  true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		sslPolicyConfigJSON := sslPolicyConfigResp.SslPolicyJSON
		if len(sslPolicyConfigJSON) > 0 {
			sslPolicy = &sslconfigs.SSLPolicy{}
			err = json.Unmarshal(sslPolicyConfigJSON, sslPolicy)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	}

	this.Data["serverType"] = server.Type
	this.Data["tlsConfig"] = maps.Map{
		"isOn":      tlsConfig.IsOn,
		"listen":    tlsConfig.Listen,
		"sslPolicy": sslPolicy,
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId   int64
	ServerType string
	Addresses  string

	SslPolicyJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.ServerTLS_LogUpdateTLSSettings, params.ServerId)

	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}

	addresses := []*serverconfigs.NetworkAddressConfig{}
	err := json.Unmarshal([]byte(params.Addresses), &addresses)
	if err != nil {
		this.Fail("端口地址解析失败：" + err.Error())
	}

	// 校验SSL
	var sslPolicyId = int64(0)
	if params.SslPolicyJSON != nil {
		sslPolicy := &sslconfigs.SSLPolicy{}
		err = json.Unmarshal(params.SslPolicyJSON, sslPolicy)
		if err != nil {
			this.ErrorPage(errors.New("解析SSL配置时发生了错误：" + err.Error()))
			return
		}

		sslPolicyId = sslPolicy.Id

		certsJSON, err := json.Marshal(sslPolicy.CertRefs)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		hstsJSON, err := json.Marshal(sslPolicy.HSTS)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		clientCACertsJSON, err := json.Marshal(sslPolicy.ClientCARefs)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		if sslPolicyId > 0 {
			_, err := this.RPC().SSLPolicyRPC().UpdateSSLPolicy(this.AdminContext(), &pb.UpdateSSLPolicyRequest{
				SslPolicyId:       sslPolicyId,
				Http2Enabled:      sslPolicy.HTTP2Enabled,
				Http3Enabled:      sslPolicy.HTTP3Enabled,
				MinVersion:        sslPolicy.MinVersion,
				SslCertsJSON:      certsJSON,
				HstsJSON:          hstsJSON,
				OcspIsOn:          sslPolicy.OCSPIsOn,
				ClientAuthType:    types.Int32(sslPolicy.ClientAuthType),
				ClientCACertsJSON: clientCACertsJSON,
				CipherSuitesIsOn:  sslPolicy.CipherSuitesIsOn,
				CipherSuites:      sslPolicy.CipherSuites,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		} else {
			resp, err := this.RPC().SSLPolicyRPC().CreateSSLPolicy(this.AdminContext(), &pb.CreateSSLPolicyRequest{
				Http2Enabled:      sslPolicy.HTTP2Enabled,
				Http3Enabled:      sslPolicy.HTTP3Enabled,
				MinVersion:        sslPolicy.MinVersion,
				SslCertsJSON:      certsJSON,
				HstsJSON:          hstsJSON,
				OcspIsOn:          sslPolicy.OCSPIsOn,
				ClientAuthType:    types.Int32(sslPolicy.ClientAuthType),
				ClientCACertsJSON: clientCACertsJSON,
				CipherSuitesIsOn:  sslPolicy.CipherSuitesIsOn,
				CipherSuites:      sslPolicy.CipherSuites,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			sslPolicyId = resp.SslPolicyId
		}
	}

	tlsConfig := &serverconfigs.TLSProtocolConfig{}
	if len(server.TlsJSON) > 0 {
		err := json.Unmarshal(server.TlsJSON, tlsConfig)
		if err != nil {
			this.ErrorPage(err)
		}
	} else {
		tlsConfig.IsOn = true
	}
	tlsConfig.Listen = addresses

	tlsConfig.SSLPolicyRef = &sslconfigs.SSLPolicyRef{
		IsOn:        true,
		SSLPolicyId: sslPolicyId,
	}

	configData, err := json.Marshal(tlsConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ServerRPC().UpdateServerTLS(this.AdminContext(), &pb.UpdateServerTLSRequest{
		ServerId: params.ServerId,
		TlsJSON:  configData,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
