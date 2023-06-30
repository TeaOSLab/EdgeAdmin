package certs

import (
	"archive/zip"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"strconv"
)

type DownloadZipAction struct {
	actionutils.ParentAction
}

func (this *DownloadZipAction) Init() {
	this.Nav("", "", "")
}

func (this *DownloadZipAction) RunGet(params struct {
	CertId int64
}) {
	defer this.CreateLogInfo(codes.SSLCert_LogDownloadSSLCertZip, params.CertId)

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

	z := zip.NewWriter(this.ResponseWriter)
	defer func() {
		_ = z.Close()
	}()

	this.AddHeader("Content-Disposition", "attachment; filename=\"cert-"+strconv.FormatInt(params.CertId, 10)+".zip\";")

	// cert
	{
		w, err := z.Create("cert.pem")
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = w.Write(certConfig.CertData)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		err = z.Flush()
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// key
	if !certConfig.IsCA {
		w, err := z.Create("key.pem")
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = w.Write(certConfig.KeyData)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		err = z.Flush()
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
}
