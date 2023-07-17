package location

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"regexp"
	"strings"
)

// IndexAction 路由规则详情
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
}

func (this *IndexAction) RunGet(params struct {
	LocationId int64
}) {
	var location = this.Data.Get("locationConfig")
	if location == nil {
		this.NotFound("location", params.LocationId)
		return
	}
	var locationConfig = location.(*serverconfigs.HTTPLocationConfig)

	this.Data["patternTypes"] = serverconfigs.AllLocationPatternTypes()

	this.Data["pattern"] = locationConfig.PatternString()
	this.Data["type"] = locationConfig.PatternType()
	this.Data["isReverse"] = locationConfig.IsReverse()
	this.Data["isCaseInsensitive"] = locationConfig.IsCaseInsensitive()
	this.Data["conds"] = locationConfig.Conds
	this.Data["domains"] = locationConfig.Domains

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	LocationId int64

	Name        string
	Pattern     string
	PatternType int
	Description string

	IsBreak           bool
	IsCaseInsensitive bool
	IsReverse         bool
	IsOn              bool

	CondsJSON []byte

	DomainsJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.HTTPLocation_LogUpdateHTTPLocation, params.LocationId)

	params.Must.
		Field("pattern", params.Pattern).
		Require("请输入路由匹配规则")

	// 校验正则
	if params.PatternType == serverconfigs.HTTPLocationPatternTypeRegexp {
		_, err := regexp.Compile(params.Pattern)
		if err != nil {
			this.Fail("正则表达式校验失败：" + err.Error())
		}
	}

	// 自动加上前缀斜杠
	if params.PatternType == serverconfigs.HTTPLocationPatternTypePrefix ||
		params.PatternType == serverconfigs.HTTPLocationPatternTypeExact {
		params.Pattern = "/" + strings.TrimLeft(params.Pattern, "/")
	}

	// 校验匹配条件
	if len(params.CondsJSON) > 0 {
		conds := &shared.HTTPRequestCondsConfig{}
		err := json.Unmarshal(params.CondsJSON, conds)
		if err != nil {
			this.Fail("匹配条件校验失败：" + err.Error())
		}

		err = conds.Init()
		if err != nil {
			this.Fail("匹配条件校验失败：" + err.Error())
		}
	}

	// 域名
	var domains = []string{}
	if len(params.DomainsJSON) > 0 {
		err := json.Unmarshal(params.DomainsJSON, &domains)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 去除可能误加的斜杠
		for index, domain := range domains {
			domains[index] = strings.TrimSuffix(domain, "/")
		}
	}

	location := &serverconfigs.HTTPLocationConfig{}
	location.SetPattern(params.Pattern, params.PatternType, params.IsCaseInsensitive, params.IsReverse)
	resultPattern := location.Pattern

	_, err := this.RPC().HTTPLocationRPC().UpdateHTTPLocation(this.AdminContext(), &pb.UpdateHTTPLocationRequest{
		LocationId:  params.LocationId,
		Name:        params.Name,
		Description: params.Description,
		Pattern:     resultPattern,
		IsBreak:     params.IsBreak,
		IsOn:        params.IsOn,
		CondsJSON:   params.CondsJSON,
		Domains:     domains,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
