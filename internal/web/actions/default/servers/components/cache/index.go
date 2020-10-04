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

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().HTTPCachePolicyRPC().CountAllEnabledHTTPCachePolicies(this.AdminContext(), &pb.CountAllEnabledHTTPCachePoliciesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	listResp, err := this.RPC().HTTPCachePolicyRPC().ListEnabledHTTPCachePolicies(this.AdminContext(), &pb.ListEnabledHTTPCachePoliciesRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cachePoliciesJSON := listResp.CachePoliciesJSON
	cachePolicies := []*serverconfigs.HTTPCachePolicy{}
	err = json.Unmarshal(cachePoliciesJSON, &cachePolicies)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["cachePolicies"] = cachePolicies

	infos := []maps.Map{}
	for _, cachePolicy := range cachePolicies {
		countServersResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithCachePolicyId(this.AdminContext(), &pb.CountAllEnabledServersWithCachePolicyIdRequest{CachePolicyId: cachePolicy.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countServers := countServersResp.Count

		infos = append(infos, maps.Map{
			"typeName":     serverconfigs.FindCachePolicyStorageName(cachePolicy.Type),
			"countServers": countServers,
		})
	}
	this.Data["infos"] = infos

	this.Show()
}
