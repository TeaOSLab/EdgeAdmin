package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

const ChinaCountryId = 1

type SelectProvincesPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectProvincesPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectProvincesPopupAction) RunGet(params struct {
	ProvinceIds string
}) {
	selectedProvinceIds := utils.SplitNumbers(params.ProvinceIds)

	provincesResp, err := this.RPC().RegionProvinceRPC().FindAllEnabledRegionProvincesWithCountryId(this.AdminContext(), &pb.FindAllEnabledRegionProvincesWithCountryIdRequest{CountryId: ChinaCountryId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	provinceMaps := []maps.Map{}
	for _, province := range provincesResp.Provinces {
		provinceMaps = append(provinceMaps, maps.Map{
			"id":        province.Id,
			"name":      province.Name,
			"isChecked": lists.ContainsInt64(selectedProvinceIds, province.Id),
		})
	}
	this.Data["provinces"] = provinceMaps

	this.Show()
}

func (this *SelectProvincesPopupAction) RunPost(params struct {
	ProvinceIds []int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	provinceMaps := []maps.Map{}
	for _, provinceId := range params.ProvinceIds {
		provinceResp, err := this.RPC().RegionProvinceRPC().FindEnabledRegionProvince(this.AdminContext(), &pb.FindEnabledRegionProvinceRequest{ProvinceId: provinceId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		province := provinceResp.Province
		if province == nil {
			continue
		}
		provinceMaps = append(provinceMaps, maps.Map{
			"id":   province.Id,
			"name": province.Name,
		})
	}
	this.Data["provinces"] = provinceMaps

	this.Success()
}
