package rewrite

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/webutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
}

func (this *CreatePopupAction) RunGet(params struct {
	WebId int64
}) {
	this.Data["webId"] = params.WebId

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	WebId          int64
	Pattern        string
	Replace        string
	Mode           string
	RedirectStatus int
	ProxyHost      string
	WithQuery      bool
	IsBreak        bool
	IsOn           bool

	Must *actions.Must
}) {
	params.Must.
		Field("pattern", params.Pattern).
		Require("请输入匹配规则").
		Expect(func() (message string, success bool) {
			_, err := regexp.Compile(params.Pattern)
			if err != nil {
				return "匹配规则错误：" + err.Error(), false
			}
			return "", true
		})

	params.Must.
		Field("replace", params.Replace).
		Require("请输入目标URL")

	// web配置
	webConfig, err := webutils.FindWebConfigWithId(this.Parent(), params.WebId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建
	createResp, err := this.RPC().HTTPRewriteRuleRPC().CreateHTTPRewriteRule(this.AdminContext(), &pb.CreateHTTPRewriteRuleRequest{
		Pattern:        params.Pattern,
		Replace:        params.Replace,
		Mode:           params.Mode,
		RedirectStatus: types.Int32(params.RedirectStatus),
		ProxyHost:      params.ProxyHost,
		WithQuery:      params.WithQuery,
		IsBreak:        params.IsBreak,
		IsOn:           params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	ref := &serverconfigs.HTTPRewriteRef{
		IsOn:          true,
		RewriteRuleId: createResp.RewriteRuleId,
	}
	webConfig.RewriteRefs = append(webConfig.RewriteRefs, ref)
	refsJSON, err := json.Marshal(webConfig.RewriteRefs)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 设置Web中的重写规则
	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebRewriteRules(this.AdminContext(), &pb.UpdateHTTPWebRewriteRulesRequest{
		WebId:            params.WebId,
		RewriteRulesJSON: refsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
