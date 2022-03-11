package certs

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type CertPopupAction struct {
	actionutils.ParentAction
}

func (this *CertPopupAction) Init() {
}

func (this *CertPopupAction) RunGet(params struct {
	CertId int64
}) {
	certResp, err := this.RPC().SSLCertRPC().FindEnabledSSLCertConfig(this.AdminContext(), &pb.FindEnabledSSLCertConfigRequest{SslCertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var certConfig = &sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(certResp.SslCertJSON, certConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var reverseCommonNames = []string{}
	for i := len(certConfig.CommonNames) - 1; i >= 0; i-- {
		reverseCommonNames = append(reverseCommonNames, certConfig.CommonNames[i])
	}

	this.Data["info"] = maps.Map{
		"id":          certConfig.Id,
		"name":        certConfig.Name,
		"description": certConfig.Description,
		"isOn":        certConfig.IsOn,
		"isAvailable": certConfig.TimeEndAt >= time.Now().Unix(),
		"commonNames": reverseCommonNames,
		"dnsNames":    certConfig.DNSNames,

		// TODO 检查是否为7天或30天内过期
		"beginTime": timeutil.FormatTime("Y-m-d H:i:s", certConfig.TimeBeginAt),
		"endTime":   timeutil.FormatTime("Y-m-d H:i:s", certConfig.TimeEndAt),

		"isCA":       certConfig.IsCA,
		"certString": string(certConfig.CertData),
		"keyString":  string(certConfig.KeyData),
	}

	// 引入的服务
	serversResp, err := this.RPC().ServerRPC().FindAllEnabledServersWithSSLCertId(this.AdminContext(), &pb.FindAllEnabledServersWithSSLCertIdRequest{SslCertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var serverMaps = []maps.Map{}
	for _, server := range serversResp.Servers {
		serverMaps = append(serverMaps, maps.Map{
			"id":   server.Id,
			"isOn": server.IsOn,
			"name": server.Name,
			"type": server.Type,
		})
	}
	this.Data["servers"] = serverMaps

	this.Show()
}
