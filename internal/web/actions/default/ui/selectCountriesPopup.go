package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type SelectCountriesPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectCountriesPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectCountriesPopupAction) RunGet(params struct {
	CountryIds string
}) {
	var selectedCountryIds = utils.SplitNumbers(params.CountryIds)

	countriesResp, err := this.RPC().RegionCountryRPC().FindAllRegionCountries(this.AdminContext(), &pb.FindAllRegionCountriesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var countryMaps = []maps.Map{}
	for _, country := range countriesResp.RegionCountries {
		countryMaps = append(countryMaps, maps.Map{
			"id":        country.Id,
			"name":      country.DisplayName,
			"letter":    strings.ToUpper(string(country.Pinyin[0][0])),
			"isChecked": lists.ContainsInt64(selectedCountryIds, country.Id),
		})
	}
	this.Data["countries"] = countryMaps

	this.Show()
}

func (this *SelectCountriesPopupAction) RunPost(params struct {
	CountryIds []int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var countryMaps = []maps.Map{}
	for _, countryId := range params.CountryIds {
		countryResp, err := this.RPC().RegionCountryRPC().FindRegionCountry(this.AdminContext(), &pb.FindRegionCountryRequest{RegionCountryId: countryId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		country := countryResp.RegionCountry
		if country == nil {
			continue
		}
		countryMaps = append(countryMaps, maps.Map{
			"id":   country.Id,
			"name": country.DisplayName,
		})
	}
	this.Data["countries"] = countryMaps

	this.Success()
}
