package server

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	adminserverutils "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/settings/server/admin-server-utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"net"
	"os"
)

type UpdateHTTPSPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateHTTPSPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateHTTPSPopupAction) RunGet(params struct{}) {
	serverConfig, err := adminserverutils.LoadServerConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["serverConfig"] = serverConfig

	// 证书
	certConfigs := []*sslconfigs.SSLCertConfig{}
	if len(serverConfig.Https.Cert) > 0 && len(serverConfig.Https.Key) > 0 {
		certData, err := os.ReadFile(Tea.Root + "/" + serverConfig.Https.Cert)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		keyData, err := os.ReadFile(Tea.Root + "/" + serverConfig.Https.Key)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		certConfig := &sslconfigs.SSLCertConfig{
			Id:       0,
			Name:     "-",
			CertData: certData,
			KeyData:  keyData,
		}
		_ = certConfig.Init(context.TODO())
		certConfig.CertData = nil
		certConfig.KeyData = nil
		certConfigs = append(certConfigs, certConfig)
	}
	this.Data["certConfigs"] = certConfigs

	this.Show()
}

func (this *UpdateHTTPSPopupAction) RunPost(params struct {
	IsOn        bool
	Listens     []string
	CertIdsJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.AdminServer_LogUpdateServerHTTPSSettings)

	if len(params.Listens) == 0 {
		this.Fail("请输入绑定地址")
	}

	serverConfig, err := adminserverutils.LoadServerConfig()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	serverConfig.Https.On = params.IsOn

	listen := []string{}
	for _, addr := range params.Listens {
		addr = utils.FormatAddress(addr)
		if len(addr) == 0 {
			continue
		}
		if _, _, err := net.SplitHostPort(addr); err != nil {
			addr += ":80"
		}
		listen = append(listen, addr)
	}
	serverConfig.Https.Listen = listen

	// 证书
	certIds := []int64{}
	err = json.Unmarshal(params.CertIdsJSON, &certIds)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if params.IsOn && len(certIds) == 0 {
		this.Fail("要启用HTTPS，需要先选择或上传一个可用的证书")
	}

	// 保存证书到本地
	if len(certIds) > 0 && certIds[0] != 0 {
		certResp, err := this.RPC().SSLCertRPC().FindEnabledSSLCertConfig(this.AdminContext(), &pb.FindEnabledSSLCertConfigRequest{
			SslCertId: certIds[0],
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(certResp.SslCertJSON) == 0 {
			this.Fail("选择的证书已失效，请换一个")
		}

		certConfig := &sslconfigs.SSLCertConfig{}
		err = json.Unmarshal(certResp.SslCertJSON, certConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		err = os.WriteFile(Tea.ConfigFile("https.key.pem"), certConfig.KeyData, 0666)
		if err != nil {
			this.Fail("保存密钥失败：" + err.Error())
		}
		err = os.WriteFile(Tea.ConfigFile("https.cert.pem"), certConfig.CertData, 0666)
		if err != nil {
			this.Fail("保存证书失败：" + err.Error())
		}

		serverConfig.Https.Key = "configs/https.key.pem"
		serverConfig.Https.Cert = "configs/https.cert.pem"
	}

	err = adminserverutils.WriteServerConfig(serverConfig)
	if err != nil {
		this.Fail("保存配置失败：" + err.Error())
	}

	this.Success()
}
