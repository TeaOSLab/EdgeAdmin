package ns

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
	UserId    int64
	Keyword   string
}) {
	if !teaconst.IsPlus {
		this.RedirectURL("/")
		return
	}

	this.Data["clusterId"] = params.ClusterId
	this.Data["userId"] = params.UserId
	this.Data["keyword"] = params.Keyword

	// 集群数量
	countClustersResp, err := this.RPC().NSClusterRPC().CountAllEnabledNSClusters(this.AdminContext(), &pb.CountAllEnabledNSClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countClusters"] = countClustersResp.Count

	// 分页
	countResp, err := this.RPC().NSDomainRPC().CountAllEnabledNSDomains(this.AdminContext(), &pb.CountAllEnabledNSDomainsRequest{
		UserId:      params.UserId,
		NsClusterId: params.ClusterId,
		Keyword:     params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	page := this.NewPage(countResp.Count)

	// 列表
	domainsResp, err := this.RPC().NSDomainRPC().ListEnabledNSDomains(this.AdminContext(), &pb.ListEnabledNSDomainsRequest{
		UserId:      params.UserId,
		NsClusterId: params.ClusterId,
		Keyword:     params.Keyword,
		Offset:      page.Offset,
		Size:        page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domainMaps := []maps.Map{}
	for _, domain := range domainsResp.NsDomains {
		// 集群信息
		var clusterMap maps.Map
		if domain.NsCluster != nil {
			clusterMap = maps.Map{
				"id":   domain.NsCluster.Id,
				"name": domain.NsCluster.Name,
			}
		}

		// 用户信息
		var userMap maps.Map
		if domain.User != nil {
			userMap = maps.Map{
				"id":       domain.User.Id,
				"username": domain.User.Username,
				"fullname": domain.User.Fullname,
			}
		}

		domainMaps = append(domainMaps, maps.Map{
			"id":      domain.Id,
			"name":    domain.Name,
			"isOn":    domain.IsOn,
			"cluster": clusterMap,
			"user":    userMap,
		})
	}
	this.Data["domains"] = domainMaps

	this.Show()
}
