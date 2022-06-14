// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type CityOptionsAction struct {
	actionutils.ParentAction
}

func (this *CityOptionsAction) RunPost(params struct{}) {
	citiesResp, err := this.RPC().RegionCityRPC().FindAllEnabledRegionCities(this.AdminContext(), &pb.FindAllEnabledRegionCitiesRequest{
		IncludeRegionProvince: true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var cityMaps = []maps.Map{}
	for _, city := range citiesResp.RegionCities {
		if city.Codes == nil {
			city.Codes = []string{}
		}

		var fullname = city.Name
		if city.RegionProvince != nil && len(city.RegionProvince.Name) > 0 && city.RegionProvince.Name != city.Name {
			fullname = city.RegionProvince.Name + " " + fullname
		}

		cityMaps = append(cityMaps, maps.Map{
			"id":       city.Id,
			"name":     city.Name,
			"fullname": fullname,
			"codes":    city.Codes,
		})
	}
	this.Data["cities"] = cityMaps

	this.Success()
}
