package userui

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
	config, err := configloaders.LoadUserUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["config"] = config

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ProductName        string
	UserSystemName     string
	ShowOpenSourceInfo bool
	ShowVersion        bool
	Version            string
	ShowFinance        bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("productName", params.ProductName).
		Require("请输入产品名称").
		Field("userSystemName", params.UserSystemName).
		Require("请输入管理员系统名称")

	config, err := configloaders.LoadUserUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	config.ProductName = params.ProductName
	config.UserSystemName = params.UserSystemName
	config.ShowOpenSourceInfo = params.ShowOpenSourceInfo
	config.ShowVersion = params.ShowVersion
	config.Version = params.Version
	config.ShowFinance = params.ShowFinance
	err = configloaders.UpdateUserUIConfig(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
