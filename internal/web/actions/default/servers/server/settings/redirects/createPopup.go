package redirects

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/url"
	"regexp"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
}) {
	this.Data["statusList"] = serverconfigs.AllHTTPRedirectStatusList()

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	BeforeURL      string
	AfterURL       string
	MatchPrefix    bool
	MatchRegexp    bool
	KeepRequestURI bool
	Status         int
	CondsJSON      []byte
	IsOn           bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("beforeURL", params.BeforeURL).
		Require("请填写跳转前的URL")

	// 校验格式
	if params.MatchRegexp {
		_, err := regexp.Compile(params.BeforeURL)
		if err != nil {
			this.Fail("跳转前URL正则表达式错误：" + err.Error())
		}
	} else {
		u, err := url.Parse(params.BeforeURL)
		if err != nil {
			this.FailField("beforeURL", "请输入正确的跳转前URL")
		}
		if (u.Scheme != "http" && u.Scheme != "https") ||
			len(u.Host) == 0 {
			this.FailField("beforeURL", "请输入正确的跳转前URL")
		}

	}

	params.Must.
		Field("afterURL", params.AfterURL).
		Require("请填写跳转后URL")

	// 校验格式
	if params.MatchRegexp {
		// 正则表达式情况下不做校验
	} else {
		u, err := url.Parse(params.AfterURL)
		if err != nil {
			this.FailField("afterURL", "请输入正确的跳转后URL")
		}
		if (u.Scheme != "http" && u.Scheme != "https") ||
			len(u.Host) == 0 {
			this.FailField("afterURL", "请输入正确的跳转后URL")
		}
	}

	params.Must.
		Field("status", params.Status).
		Gte(0, "请选择正确的跳转状态码")

	// 校验匹配条件
	var conds *shared.HTTPRequestCondsConfig
	if len(params.CondsJSON) > 0 {
		conds = &shared.HTTPRequestCondsConfig{}
		err := json.Unmarshal(params.CondsJSON, conds)
		if err != nil {
			this.Fail("匹配条件校验失败：" + err.Error())
		}

		err = conds.Init()
		if err != nil {
			this.Fail("匹配条件校验失败：" + err.Error())
		}
	}

	this.Data["redirect"] = maps.Map{
		"status":         params.Status,
		"beforeURL":      params.BeforeURL,
		"afterURL":       params.AfterURL,
		"matchPrefix":    params.MatchPrefix,
		"matchRegexp":    params.MatchRegexp,
		"keepRequestURI": params.KeepRequestURI,
		"conds":          conds,
		"isOn":           params.IsOn,
	}

	this.Success()
}
