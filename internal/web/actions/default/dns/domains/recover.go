package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type RecoverAction struct {
	actionutils.ParentAction
}

func (this *RecoverAction) RunPost(params struct {
	DomainId int64
}) {
	// 记录日志
	defer this.CreateLog(oplogs.LevelInfo, "从DNS服务商中恢复域名 %d", params.DomainId)

	// 执行恢复
	_, err := this.RPC().DNSDomainRPC().RecoverDNSDomain(this.AdminContext(), &pb.RecoverDNSDomainRequest{DnsDomainId: params.DomainId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
