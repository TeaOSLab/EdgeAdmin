package server

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
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
	config, err := configloaders.LoadUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["config"] = config

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ProductName        string
	AdminSystemName    string
	ShowOpenSourceInfo bool
	ShowVersion        bool
	Version            string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("productName", params.ProductName).
		Require("请输入产品名称").
		Field("adminSystemName", params.AdminSystemName).
		Require("请输入管理员系统名称")

	config, err := configloaders.LoadUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	config.ProductName = params.ProductName
	config.AdminSystemName = params.AdminSystemName
	config.ShowOpenSourceInfo = params.ShowOpenSourceInfo
	config.ShowVersion = params.ShowVersion
	config.Version = params.Version
	err = configloaders.UpdateUIConfig(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
