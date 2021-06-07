package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	Keyword string
}) {
	this.Data["keyword"] = params.Keyword

	countResp, err := this.RPC().HTTPCachePolicyRPC().CountAllEnabledHTTPCachePolicies(this.AdminContext(), &pb.CountAllEnabledHTTPCachePoliciesRequest{
		Keyword: params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	listResp, err := this.RPC().HTTPCachePolicyRPC().ListEnabledHTTPCachePolicies(this.AdminContext(), &pb.ListEnabledHTTPCachePoliciesRequest{
		Keyword: params.Keyword,
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cachePoliciesJSON := listResp.HttpCachePoliciesJSON
	cachePolicies := []*serverconfigs.HTTPCachePolicy{}
	err = json.Unmarshal(cachePoliciesJSON, &cachePolicies)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["cachePolicies"] = cachePolicies

	infos := []maps.Map{}
	for _, cachePolicy := range cachePolicies {
		countClustersResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClustersWithHTTPCachePolicyId(this.AdminContext(), &pb.CountAllEnabledNodeClustersWithHTTPCachePolicyIdRequest{HttpCachePolicyId: cachePolicy.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countClusters := countClustersResp.Count

		infos = append(infos, maps.Map{
			"typeName":      serverconfigs.FindCachePolicyStorageName(cachePolicy.Type),
			"countClusters": countClusters,
		})
	}
	this.Data["infos"] = infos

	this.Show()
}
