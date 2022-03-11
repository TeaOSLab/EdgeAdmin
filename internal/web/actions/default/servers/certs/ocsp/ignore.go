// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ocsp

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type IgnoreAction struct {
	actionutils.ParentAction
}

func (this *IgnoreAction) RunPost(params struct {
	CertIds []int64
}) {
	defer this.CreateLogInfo("忽略一组证书的OCSP状态")

	_, err := this.RPC().SSLCertRPC().IgnoreSSLCertsWithOCSPError(this.AdminContext(), &pb.IgnoreSSLCertsWithOCSPErrorRequest{SslCertIds: params.CertIds})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
