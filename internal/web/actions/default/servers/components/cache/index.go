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
	ClusterId   int64
	Keyword     string
	StorageType string
}) {
	this.Data["keyword"] = params.Keyword
	this.Data["clusterId"] = params.ClusterId
	this.Data["storageType"] = params.StorageType

	countResp, err := this.RPC().HTTPCachePolicyRPC().CountAllEnabledHTTPCachePolicies(this.AdminContext(), &pb.CountAllEnabledHTTPCachePoliciesRequest{
		NodeClusterId: params.ClusterId,
		Keyword:       params.Keyword,
		Type:          params.StorageType,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	listResp, err := this.RPC().HTTPCachePolicyRPC().ListEnabledHTTPCachePolicies(this.AdminContext(), &pb.ListEnabledHTTPCachePoliciesRequest{
		Keyword:       params.Keyword,
		NodeClusterId: params.ClusterId,
		Type:          params.StorageType,
		Offset:        page.Offset,
		Size:          page.Size,
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

	// 所有的存储类型
	this.Data["storageTypes"] = serverconfigs.AllCachePolicyStorageTypes

	this.Show()
}
