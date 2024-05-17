// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"strings"
)

type TestRegexpAction struct {
	actionutils.ParentAction
}

func (this *TestRegexpAction) RunPost(params struct {
	Regexp            string
	IsCaseInsensitive bool
	Body              string
}) {
	var exp = params.Regexp
	if params.IsCaseInsensitive && !strings.HasPrefix(params.Regexp, "(?i)") {
		exp = "(?i)" + exp
	}
	reg, err := regexp.Compile(exp)
	if err != nil {
		this.Data["result"] = maps.Map{
			"isOk":    false,
			"message": "解析正则出错：" + err.Error(),
		}
		this.Success()
	}

	if reg.MatchString(params.Body) {
		this.Data["result"] = maps.Map{
			"isOk":    true,
			"message": "匹配成功",
		}
		this.Success()
	}

	this.Data["result"] = maps.Map{
		"isOk":    false,
		"message": "匹配失败",
	}

	this.Success()
}
