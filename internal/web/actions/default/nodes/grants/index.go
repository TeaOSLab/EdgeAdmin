package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/grants/grantutils"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "grant", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().NodeGrantRPC().CountAllEnabledNodeGrants(this.AdminContext(), &pb.CountAllEnabledNodeGrantsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	page := this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	grantsResp, err := this.RPC().NodeGrantRPC().ListEnabledNodeGrants(this.AdminContext(), &pb.ListEnabledNodeGrantsRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	grantMaps := []maps.Map{}
	for _, grant := range grantsResp.Grants {
		grantMaps = append(grantMaps, maps.Map{
			"id":   grant.Id,
			"name": grant.Name,
			"method": maps.Map{
				"type": grant.Method,
				"name": grantutils.FindGrantMethodName(grant.Method),
			},
		})
	}
	this.Data["grants"] = grantMaps

	this.Show()
}
