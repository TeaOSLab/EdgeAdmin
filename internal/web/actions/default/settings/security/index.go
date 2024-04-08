package security

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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

func (this *IndexAction) RunGet(params struct {
	ShowAll bool
}) {
	this.Data["showAll"] = params.ShowAll

	config, err := configloaders.LoadSecurityConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if config.AllowIPs == nil {
		config.AllowIPs = []string{}
	}

	// 国家和地区
	var countryMaps = []maps.Map{}
	for _, countryId := range config.AllowCountryIds {
		countryResp, err := this.RPC().RegionCountryRPC().FindRegionCountry(this.AdminContext(), &pb.FindRegionCountryRequest{RegionCountryId: countryId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var country = countryResp.RegionCountry
		if country != nil {
			countryMaps = append(countryMaps, maps.Map{
				"id":   country.Id,
				"name": country.DisplayName,
			})
		}
	}
	this.Data["countries"] = countryMaps

	// 省份
	var provinceMaps = []maps.Map{}
	for _, provinceId := range config.AllowProvinceIds {
		provinceResp, err := this.RPC().RegionProvinceRPC().FindRegionProvince(this.AdminContext(), &pb.FindRegionProvinceRequest{RegionProvinceId: provinceId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var province = provinceResp.RegionProvince
		if province != nil {
			provinceMaps = append(provinceMaps, maps.Map{
				"id":   province.Id,
				"name": province.DisplayName,
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

	ClientIPHeaderNames string
	ClientIPHeaderOnly  bool

	DenySearchEngines bool
	DenySpiders       bool

	CheckClientFingerprint bool
	CheckClientRegion      bool

	DomainsJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.AdminSecurity_LogUpdateSecuritySettings)

	config, err := configloaders.LoadSecurityConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 框架
	config.Frame = params.Frame

	// 国家和地区
	var countryIds = []int64{}
	if len(params.CountryIdsJSON) > 0 {
		err = json.Unmarshal(params.CountryIdsJSON, &countryIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	config.AllowCountryIds = countryIds

	// 省份
	var provinceIds = []int64{}
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

	// 客户端IP获取方式
	config.ClientIPHeaderNames = params.ClientIPHeaderNames
	config.ClientIPHeaderOnly = params.ClientIPHeaderOnly

	// 禁止搜索引擎和爬虫
	config.DenySearchEngines = params.DenySearchEngines
	config.DenySpiders = params.DenySpiders

	// 允许的域名
	var domains = []string{}
	if len(params.DomainsJSON) > 0 {
		err = json.Unmarshal(params.DomainsJSON, &domains)
		if err != nil {
			this.Fail("解析允许访问的域名失败：" + err.Error())
		}
	}
	config.AllowDomains = domains

	// 允许记住登录
	config.AllowRememberLogin = params.AllowRememberLogin

	// Cookie检查
	config.CheckClientFingerprint = params.CheckClientFingerprint
	config.CheckClientRegion = params.CheckClientRegion

	err = configloaders.UpdateSecurityConfig(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
