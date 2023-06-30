package locations

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"regexp"
	"strings"
)

// CreateAction 创建路由规则
type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "setting", "create")
	this.SecondMenu("locations")
}

func (this *CreateAction) RunGet(params struct {
	ServerId int64
	ParentId int64 // 父节点
}) {
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["patternTypes"] = serverconfigs.AllLocationPatternTypes()

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	WebId int64

	Name        string
	Pattern     string
	PatternType int
	Description string

	IsBreak           bool
	IsCaseInsensitive bool
	IsReverse         bool
	CondsJSON         []byte

	DomainsJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.HTTPLocation_LogCreateHTTPLocation, params.Pattern)

	params.Must.
		Field("pattern", params.Pattern).
		Require("请输入路径匹配规则")

	// 校验正则
	if params.PatternType == serverconfigs.HTTPLocationPatternTypeRegexp {
		_, err := regexp.Compile(params.Pattern)
		if err != nil {
			this.Fail("正则表达式校验失败：" + err.Error())
		}
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

	// 自动加上前缀斜杠
	if params.PatternType == serverconfigs.HTTPLocationPatternTypePrefix ||
		params.PatternType == serverconfigs.HTTPLocationPatternTypeExact {
		params.Pattern = "/" + strings.TrimLeft(params.Pattern, "/")
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

	locationResp, err := this.RPC().HTTPLocationRPC().CreateHTTPLocation(this.AdminContext(), &pb.CreateHTTPLocationRequest{
		ParentId:    0, // TODO 需要实现
		Name:        params.Name,
		Description: params.Description,
		Pattern:     resultPattern,
		IsBreak:     params.IsBreak,
		CondsJSON:   params.CondsJSON,
		Domains:     domains,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	locationId := locationResp.LocationId

	// Web中Location
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithId(this.AdminContext(), params.WebId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	// TODO 支持Location嵌套
	webConfig.LocationRefs = append(webConfig.LocationRefs, &serverconfigs.HTTPLocationRef{
		IsOn:       true,
		LocationId: locationId,
		Children:   nil,
	})

	refJSON, err := json.Marshal(webConfig.LocationRefs)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebLocations(this.AdminContext(), &pb.UpdateHTTPWebLocationsRequest{
		HttpWebId:     params.WebId,
		LocationsJSON: refJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
