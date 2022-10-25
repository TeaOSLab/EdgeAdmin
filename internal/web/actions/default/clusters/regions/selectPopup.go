package regions

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectPopupAction) RunGet(params struct{}) {
	regionsResp, err := this.RPC().NodeRegionRPC().FindAllAvailableNodeRegions(this.AdminContext(), &pb.FindAllAvailableNodeRegionsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	regionMaps := []maps.Map{}
	for _, region := range regionsResp.NodeRegions {
		regionMaps = append(regionMaps, maps.Map{
			"id":   region.Id,
			"name": region.Name,
		})
	}
	this.Data["regions"] = regionMaps

	this.Show()
}

func (this *SelectPopupAction) RunPost(params struct {
	RegionId int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	if params.RegionId <= 0 {
		this.Data["region"] = nil
		this.Success()
		return
	}

	regionResp, err := this.RPC().NodeRegionRPC().FindEnabledNodeRegion(this.AdminContext(), &pb.FindEnabledNodeRegionRequest{NodeRegionId: params.RegionId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var region = regionResp.NodeRegion
	if region == nil {
		this.NotFound("nodeRegion", params.RegionId)
		return
	}

	this.Data["region"] = maps.Map{
		"id":   region.Id,
		"name": region.Name,
	}

	this.Success()
}
