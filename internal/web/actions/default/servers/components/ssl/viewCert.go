package ssl

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
)

type ViewCertAction struct {
	actionutils.ParentAction
}

func (this *ViewCertAction) Init() {
	this.Nav("", "", "")
}

func (this *ViewCertAction) RunGet(params struct {
	CertId int64
}) {
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
	this.Write(certConfig.CertData)
}
