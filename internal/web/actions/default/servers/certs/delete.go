package certs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	CertId int64
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "删除SSL证书 %d", params.CertId)

	// 是否正在被使用
	countResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithSSLCertId(this.AdminContext(), &pb.CountAllEnabledServersWithSSLCertIdRequest{SslCertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count > 0 {
		this.Fail("此证书正在被某些服务引用，请先修改服务后再删除。")
	}

	_, err = this.RPC().SSLCertRPC().DeleteSSLCert(this.AdminContext(), &pb.DeleteSSLCertRequest{SslCertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
