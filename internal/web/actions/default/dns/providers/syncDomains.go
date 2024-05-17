// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package providers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type SyncDomainsAction struct {
	actionutils.ParentAction
}

func (this *SyncDomainsAction) RunPost(params struct {
	ProviderId int64
}) {
	resp, err := this.RPC().DNSDomainRPC().SyncDNSDomainsFromProvider(this.AdminContext(), &pb.SyncDNSDomainsFromProviderRequest{DnsProviderId: params.ProviderId})
	if err != nil {
		this.Fail("更新域名失败：" + err.Error())
	}

	this.Data["hasChanges"] = resp.HasChanges

	this.Success()
}
