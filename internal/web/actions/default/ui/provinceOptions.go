// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/regionconfigs"
	"github.com/iwind/TeaGo/maps"
)

type ProvinceOptionsAction struct {
	actionutils.ParentAction
}

func (this *ProvinceOptionsAction) RunPost(params struct{}) {
	provincesResp, err := this.RPC().RegionProvinceRPC().FindAllRegionProvincesWithRegionCountryId(this.AdminContext(), &pb.FindAllRegionProvincesWithRegionCountryIdRequest{RegionCountryId: regionconfigs.RegionChinaId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var provinceMaps = []maps.Map{}
	for _, province := range provincesResp.RegionProvinces {
		if province.Codes == nil {
			province.Codes = []string{}
		}
		provinceMaps = append(provinceMaps, maps.Map{
			"id":    province.Id,
			"name":  province.DisplayName,
			"codes": province.Codes,
		})
	}
	this.Data["provinces"] = provinceMaps

	this.Success()
}
