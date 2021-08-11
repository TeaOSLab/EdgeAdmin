package security

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	config, err := configloaders.LoadSecurityConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if config.AllowIPs == nil {
		config.AllowIPs = []string{}
	}

	// 国家和地区
	countryMaps := []maps.Map{}
	for _, countryId := range config.AllowCountryIds {
		countryResp, err := this.RPC().RegionCountryRPC().FindEnabledRegionCountry(this.AdminContext(), &pb.FindEnabledRegionCountryRequest{CountryId: countryId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		country := countryResp.Country
		if country != nil {
			countryMaps = append(countryMaps, maps.Map{
				"id":   country.Id,
				"name": country.Name,
			})
		}
	}
	this.Data["countries"] = countryMaps

	// 省份
	provinceMaps := []maps.Map{}
	for _, provinceId := range config.AllowProvinceIds {
		provinceResp, err := this.RPC().RegionProvinceRPC().FindEnabledRegionProvince(this.AdminContext(), &pb.FindEnabledRegionProvinceRequest{ProvinceId: provinceId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		province := provinceResp.Province
		if province != nil {
			provinceMaps = append(provinceMaps, maps.Map{
				"id":   province.Id,
				"name": province.Name,
			})
		}
	}
	this.Data["provinces"] = provinceMaps

	this.Data["config"] = config
	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	Frame              string
	CountryIdsJSON     []byte
	ProvinceIdsJSON    []byte
	AllowLocal         bool
	AllowIPs           []string
	AllowRememberLogin bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改管理界面安全设置")

	config, err := configloaders.LoadSecurityConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 框架
	config.Frame = params.Frame

	// 国家和地区
	countryIds := []int64{}
	if len(params.CountryIdsJSON) > 0 {
		err = json.Unmarshal(params.CountryIdsJSON, &countryIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	config.AllowCountryIds = countryIds

	// 省份
	provinceIds := []int64{}
	if len(params.ProvinceIdsJSON) > 0 {
		err = json.Unmarshal(params.ProvinceIdsJSON, &provinceIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	config.AllowProvinceIds = provinceIds

	// 允许的IP
	if len(params.AllowIPs) > 0 {
		for _, ip := range params.AllowIPs {
			_, err := shared.ParseIPRange(ip)
			if err != nil {
				this.Fail("允许访问的IP '" + ip + "' 格式错误：" + err.Error())
			}
		}
		config.AllowIPs = params.AllowIPs
	} else {
		config.AllowIPs = []string{}
	}

	// 允许本地
	config.AllowLocal = params.AllowLocal

	// 允许记住登录
	config.AllowRememberLogin = params.AllowRememberLogin

	err = configloaders.UpdateSecurityConfig(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
