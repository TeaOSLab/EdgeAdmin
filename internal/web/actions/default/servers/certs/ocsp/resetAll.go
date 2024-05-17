// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package ocsp

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type ResetAllAction struct {
	actionutils.ParentAction
}

func (this *ResetAllAction) RunPost(params struct{}) {
	defer this.CreateLogInfo(codes.SSLCert_LogOCSPResetAllOCSPStatus)

	_, err := this.RPC().SSLCertRPC().ResetAllSSLCertsWithOCSPError(this.AdminContext(), &pb.ResetAllSSLCertsWithOCSPErrorRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
