package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) Init() {
	this.Nav("", "delete", "index")
	this.SecondMenu("nodes")
}

func (this *DeleteAction) RunGet(params struct{}) {
	this.Show()
}

func (this *DeleteAction) RunPost(params struct {
	ClusterId int64
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "删除域名服务集群 %d", params.ClusterId)

	// TODO 如果有用户在使用此集群，就不能删除

	// 删除
	_, err := this.RPC().NSClusterRPC().DeleteNSCluster(this.AdminContext(), &pb.DeleteNSCluster{NsClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
