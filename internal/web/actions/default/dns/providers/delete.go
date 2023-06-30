package providers

import (	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	ProviderId int64
}) {
	// TODO 检查权限

	// 记录日志
	defer this.CreateLogInfo(codes.DNSProvider_LogDeleteDNSProvider, params.ProviderId)

	// 检查是否被集群使用
	countClustersResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClustersWithDNSProviderId(this.AdminContext(), &pb.CountAllEnabledNodeClustersWithDNSProviderIdRequest{DnsProviderId: params.ProviderId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countClustersResp.Count > 0 {
		this.Fail("当前DNS服务商账号正在被" + numberutils.FormatInt64(countClustersResp.Count) + "个集群所使用，所以不能删除。请修改集群设置后再操作。")
	}

	// 判断是否被ACME任务使用
	countACMETasksResp, err := this.RPC().ACMETaskRPC().CountEnabledACMETasksWithDNSProviderId(this.AdminContext(), &pb.CountEnabledACMETasksWithDNSProviderIdRequest{DnsProviderId: params.ProviderId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countACMETasksResp.Count > 0 {
		this.Fail("当前DNS服务商账号正在被" + numberutils.FormatInt64(countACMETasksResp.Count) + "个ACME证书申请任务使用，所以不能删除。请修改ACME证书申请任务后再操作。")
	}

	// 执行删除
	_, err = this.RPC().DNSProviderRPC().DeleteDNSProvider(this.AdminContext(), &pb.DeleteDNSProviderRequest{DnsProviderId: params.ProviderId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
