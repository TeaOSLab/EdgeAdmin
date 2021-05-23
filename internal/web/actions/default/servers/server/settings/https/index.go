package https

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("https")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	httpsConfig := &serverconfigs.HTTPSProtocolConfig{}
	if len(server.HttpsJSON) > 0 {
		err := json.Unmarshal(server.HttpsJSON, httpsConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		httpsConfig.IsOn = true
	}

	var sslPolicy *sslconfigs.SSLPolicy
	if httpsConfig.SSLPolicyRef != nil && httpsConfig.SSLPolicyRef.SSLPolicyId > 0 {
		sslPolicyConfigResp, err := this.RPC().SSLPolicyRPC().FindEnabledSSLPolicyConfig(this.AdminContext(), &pb.FindEnabledSSLPolicyConfigRequest{SslPolicyId: httpsConfig.SSLPolicyRef.SSLPolicyId})
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
	this.Data["httpsConfig"] = maps.Map{
		"isOn":      httpsConfig.IsOn,
		"addresses": httpsConfig.Listen,
		"sslPolicy": sslPolicy,
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId  int64
	IsOn      bool
	Addresses string

	SslPolicyJSON []byte

	Must *actions.Must
}) {
	// 记录日志
	defer this.CreateLog(oplogs.LevelInfo, "修改服务 %d 的HTTPS设置", params.ServerId)

	addresses := []*serverconfigs.NetworkAddressConfig{}
	err := json.Unmarshal([]byte(params.Addresses), &addresses)
	if err != nil {
		this.Fail("端口地址解析失败：" + err.Error())
	}

	// 检查端口地址是否正确
	for _, addr := range addresses {
		err = addr.Init()
		if err != nil {
			this.Fail("绑定端口校验失败：" + err.Error())
		}

		if regexp.MustCompile(`^\d+$`).MatchString(addr.PortRange) {
			port := types.Int(addr.PortRange)
			if port > 65535 {
				this.Fail("绑定的端口地址不能大于65535")
			}
			if port == 80 {
				this.Fail("端口80通常是HTTP的端口，不能用在HTTPS上")
			}
		}
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
				MinVersion:        sslPolicy.MinVersion,
				SslCertsJSON:      certsJSON,
				HstsJSON:          hstsJSON,
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
				MinVersion:        sslPolicy.MinVersion,
				SslCertsJSON:      certsJSON,
				HstsJSON:          hstsJSON,
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

	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	httpsConfig := &serverconfigs.HTTPSProtocolConfig{}
	if len(server.HttpsJSON) > 0 {
		err = json.Unmarshal(server.HttpsJSON, httpsConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	httpsConfig.SSLPolicyRef = &sslconfigs.SSLPolicyRef{
		IsOn:        true,
		SSLPolicyId: sslPolicyId,
	}
	httpsConfig.IsOn = params.IsOn
	httpsConfig.Listen = addresses
	configData, err := json.Marshal(httpsConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ServerRPC().UpdateServerHTTPS(this.AdminContext(), &pb.UpdateServerHTTPSRequest{
		ServerId:  params.ServerId,
		HttpsJSON: configData,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
