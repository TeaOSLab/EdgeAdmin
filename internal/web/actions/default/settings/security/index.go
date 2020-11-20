package security

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/securitymanager"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	config, err := securitymanager.LoadSecurityConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["config"] = config
	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	Frame string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改管理界面安全设置")

	config, err := securitymanager.LoadSecurityConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	config.Frame = params.Frame
	err = securitymanager.UpdateSecurityConfig(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
