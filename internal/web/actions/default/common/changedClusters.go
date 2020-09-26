package common

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"time"
)

// 检查变更的集群列表
type ChangedClustersAction struct {
	actionutils.ParentAction
}

func (this *ChangedClustersAction) Init() {
	this.Nav("", "", "")
}

func (this *ChangedClustersAction) RunGet(params struct {
	IsNotifying bool
}) {
	timeout := time.NewTimer(55 * time.Second) // 比客户端提前结束，避免在客户端产生一个请求错误

	this.Data["clusters"] = []interface{}{}

Loop:
	for {
		select {
		case <-this.Request.Context().Done():
			break Loop
		case <-timeout.C:
			break Loop
		default:
			// 继续
		}

		resp, err := this.RPC().NodeClusterRPC().FindAllChangedNodeClusters(this.AdminContext(), &pb.FindAllChangedNodeClustersRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		result := []maps.Map{}
		for _, cluster := range resp.Clusters {
			result = append(result, maps.Map{
				"id":   cluster.Id,
				"name": cluster.Name,
			})
		}

		// 从提醒到提醒消失
		if len(result) == 0 && params.IsNotifying {
			break
		}

		this.Data["clusters"] = result
		if len(result) > 0 {
			break
		}

		time.Sleep(1 * time.Second)
	}

	this.Success()
}
