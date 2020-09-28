package rewrite

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/webutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
}

func (this *UpdatePopupAction) RunGet(params struct {
	WebId         int64
	RewriteRuleId int64
}) {
	this.Data["webId"] = params.WebId

	webConfig, err := webutils.FindWebConfigWithId(this.Parent(), params.WebId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	isFound := false
	for _, rewriteRule := range webConfig.RewriteRules {
		if rewriteRule.Id == params.RewriteRuleId {
			this.Data["rewriteRule"] = rewriteRule
			isFound = true
			break
		}
	}

	if !isFound {
		this.WriteString("找不到要修改的重写规则")
		return
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	WebId          int64
	RewriteRuleId  int64
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

	// 修改
	_, err := this.RPC().HTTPRewriteRuleRPC().UpdateHTTPRewriteRule(this.AdminContext(), &pb.UpdateHTTPRewriteRuleRequest{
		RewriteRuleId:  params.RewriteRuleId,
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

	this.Success()
}
