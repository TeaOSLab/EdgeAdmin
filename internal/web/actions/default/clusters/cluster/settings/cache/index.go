package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("cache")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	cluster, err := dao.SharedNodeClusterDAO.FindEnabledNodeCluster(this.AdminContext(), params.ClusterId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if cluster == nil {
		this.NotFound("nodeCluster", params.ClusterId)
		return
	}

	// 缓存设置
	this.Data["cachePolicy"] = nil
	if cluster.HttpCachePolicyId > 0 {
		cachePolicy, err := dao.SharedHTTPCachePolicyDAO.FindEnabledHTTPCachePolicy(this.AdminContext(), cluster.HttpCachePolicyId)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if cachePolicy != nil {
			this.Data["cachePolicy"] = maps.Map{
				"id":   cachePolicy.Id,
				"name": cachePolicy.Name,
				"isOn": cachePolicy.IsOn,
			}
		}
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId     int64
	CachePolicyId int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.ServerCache_LogUpdateClusterCachePolicy, params.ClusterId, params.CachePolicyId)

	if params.CachePolicyId <= 0 {
		this.Fail("请选择缓存策略")
	}

	_, err := this.RPC().NodeClusterRPC().UpdateNodeClusterHTTPCachePolicyId(this.AdminContext(), &pb.UpdateNodeClusterHTTPCachePolicyIdRequest{
		NodeClusterId:     params.ClusterId,
		HttpCachePolicyId: params.CachePolicyId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
