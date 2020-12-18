package certs

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"strconv"
)

type DownloadCertAction struct {
	actionutils.ParentAction
}

func (this *DownloadCertAction) Init() {
	this.Nav("", "", "")
}

func (this *DownloadCertAction) RunGet(params struct {
	CertId int64
}) {
	defer this.CreateLogInfo("下载SSL证书 %d", params.CertId)

	certResp, err := this.RPC().SSLCertRPC().FindEnabledSSLCertConfig(this.AdminContext(), &pb.FindEnabledSSLCertConfigRequest{SslCertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	certConfig := &sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(certResp.SslCertJSON, certConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.AddHeader("Content-Disposition", "attachment; filename=\"cert-"+strconv.FormatInt(params.CertId, 10)+".pem\";")
	this.Write(certConfig.CertData)
}
