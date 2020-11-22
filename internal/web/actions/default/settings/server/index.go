package server

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	this.Data["serverIsChanged"] = serverConfigIsChanged

	serverConfig, err := loadServerConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["serverConfig"] = serverConfig

	this.Show()
}
