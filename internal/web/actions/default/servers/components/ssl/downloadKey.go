package ssl

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"strconv"
)

type DownloadKeyAction struct {
	actionutils.ParentAction
}

func (this *DownloadKeyAction) Init() {
	this.Nav("", "", "")
}

func (this *DownloadKeyAction) RunGet(params struct {
	CertId int64
}) {
	defer this.CreateLogInfo("下载SSL密钥 %d", params.CertId)

	certResp, err := this.RPC().SSLCertRPC().FindEnabledSSLCertConfig(this.AdminContext(), &pb.FindEnabledSSLCertConfigRequest{CertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	certConfig := &sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(certResp.CertJSON, certConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.AddHeader("Content-Disposition", "attachment; filename=\"key-"+strconv.FormatInt(params.CertId, 10)+".pem\";")
	this.Write(certConfig.KeyData)
}
