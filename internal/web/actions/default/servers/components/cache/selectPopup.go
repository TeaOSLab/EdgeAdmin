package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectPopupAction) RunGet(params struct{}) {
	countResp, err := this.RPC().HTTPCachePolicyRPC().CountAllEnabledHTTPCachePolicies(this.AdminContext(), &pb.CountAllEnabledHTTPCachePoliciesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)

	this.Data["page"] = page.AsHTML()

	cachePoliciesResp, err := this.RPC().HTTPCachePolicyRPC().ListEnabledHTTPCachePolicies(this.AdminContext(), &pb.ListEnabledHTTPCachePoliciesRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cachePolicies := []*serverconfigs.HTTPCachePolicy{}
	if len(cachePoliciesResp.HttpCachePoliciesJSON) > 0 {
		err = json.Unmarshal(cachePoliciesResp.HttpCachePoliciesJSON, &cachePolicies)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	policyMaps := []maps.Map{}
	for _, cachePolicy := range cachePolicies {
		policyMaps = append(policyMaps, maps.Map{
			"id":          cachePolicy.Id,
			"name":        cachePolicy.Name,
			"description": cachePolicy.Description,
			"isOn":        cachePolicy.IsOn,
		})
	}

	this.Data["cachePolicies"] = policyMaps

	this.Show()
}
