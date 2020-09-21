package locations

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/webutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"regexp"
	"strings"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "setting", "create")
	this.SecondMenu("locations")
}

func (this *CreateAction) RunGet(params struct {
	ServerId int64
}) {
	webConfig, err := webutils.FindWebConfigWithServerId(this.Parent(), params.ServerId)
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

	Must *actions.Must
}) {
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

	// 自动加上前缀斜杠
	if params.PatternType == serverconfigs.HTTPLocationPatternTypePrefix ||
		params.PatternType == serverconfigs.HTTPLocationPatternTypeExact {
		params.Pattern = "/" + strings.TrimLeft(params.Pattern, "/")
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
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	locationId := locationResp.LocationId

	// Web中Location
	webConfig, err := webutils.FindWebConfigWithId(this.Parent(), params.WebId)
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
		WebId:         params.WebId,
		LocationsJSON: refJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
