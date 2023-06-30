package regions

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	RegionId int64
}) {
	regionResp, err := this.RPC().NodeRegionRPC().FindEnabledNodeRegion(this.AdminContext(), &pb.FindEnabledNodeRegionRequest{NodeRegionId: params.RegionId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	region := regionResp.NodeRegion
	if region == nil {
		this.NotFound("nodeRegion", params.RegionId)
		return
	}

	this.Data["region"] = maps.Map{
		"id":          region.Id,
		"isOn":        region.IsOn,
		"name":        region.Name,
		"description": region.Description,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	RegionId int64

	Name        string
	Description string
	IsOn        bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.NodeRegion_LogUpdateNodeRegion, params.RegionId)

	params.Must.
		Field("name", params.Name).
		Require("请输入区域名称")

	_, err := this.RPC().NodeRegionRPC().UpdateNodeRegion(this.AdminContext(), &pb.UpdateNodeRegionRequest{
		NodeRegionId: params.RegionId,
		Name:         params.Name,
		Description:  params.Description,
		IsOn:         params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
