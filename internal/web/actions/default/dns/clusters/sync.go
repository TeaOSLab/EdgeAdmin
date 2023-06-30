package clusters

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type SyncAction struct {
	actionutils.ParentAction
}

func (this *SyncAction) RunPost(params struct {
	ClusterId int64
}) {
	// 记录日志
	defer this.CreateLogInfo(codes.DNS_LogSyncCluster, params.ClusterId)

	dnsInfoResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterDNS(this.AdminContext(), &pb.FindEnabledNodeClusterDNSRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domain := dnsInfoResp.Domain
	if domain == nil || domain.Id <= 0 {
		this.Fail("此集群尚未设置域名")
	}

	syncResp, err := this.RPC().DNSDomainRPC().SyncDNSDomainData(this.AdminContext(), &pb.SyncDNSDomainDataRequest{
		DnsDomainId:   domain.Id,
		NodeClusterId: params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if syncResp.ShouldFix {
		this.Fail("请先修改当前页面中红色标记的问题")
	}

	if !syncResp.IsOk {
		this.Fail(syncResp.Error)
	}

	this.Success()
}
