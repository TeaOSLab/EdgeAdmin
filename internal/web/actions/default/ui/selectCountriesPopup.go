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
	selectedCountryIds := utils.SplitNumbers(params.CountryIds)

	countriesResp, err := this.RPC().RegionCountryRPC().FindAllEnabledRegionCountries(this.AdminContext(), &pb.FindAllEnabledRegionCountriesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	countryMaps := []maps.Map{}
	for _, country := range countriesResp.Countries {
		countryMaps = append(countryMaps, maps.Map{
			"id":        country.Id,
			"name":      country.Name,
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
	countryMaps := []maps.Map{}
	for _, countryId := range params.CountryIds {
		countryResp, err := this.RPC().RegionCountryRPC().FindEnabledRegionCountry(this.AdminContext(), &pb.FindEnabledRegionCountryRequest{CountryId: countryId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		country := countryResp.Country
		if country == nil {
			continue
		}
		countryMaps = append(countryMaps, maps.Map{
			"id":   country.Id,
			"name": country.Name,
		})
	}
	this.Data["countries"] = countryMaps

	this.Success()
}
