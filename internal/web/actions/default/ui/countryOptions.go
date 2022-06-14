// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type CountryOptionsAction struct {
	actionutils.ParentAction
}

func (this *CountryOptionsAction) RunPost(params struct{}) {
	countriesResp, err := this.RPC().RegionCountryRPC().FindAllEnabledRegionCountries(this.AdminContext(), &pb.FindAllEnabledRegionCountriesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var countryMaps = []maps.Map{}
	for _, country := range countriesResp.RegionCountries {
		if country.Codes == nil {
			country.Codes = []string{}
		}

		var letter = ""
		if len(country.Pinyin) > 0 && len(country.Pinyin) > 0 && len(country.Pinyin[0]) > 0 {
			letter = strings.ToUpper(country.Pinyin[0][:1])
		}

		countryMaps = append(countryMaps, maps.Map{
			"id":       country.Id,
			"name":     country.Name,
			"fullname": letter + " " + country.Name,
			"codes":    country.Codes,
		})
	}
	this.Data["countries"] = countryMaps

	this.Success()
}
