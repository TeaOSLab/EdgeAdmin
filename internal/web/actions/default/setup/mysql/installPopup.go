// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package mysql

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/setup/mysql/mysqlinstallers"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/setup/mysql/mysqlinstallers/utils"
)

type InstallPopupAction struct {
	actionutils.ParentAction
}

func (this *InstallPopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *InstallPopupAction) RunPost(params struct{}) {
	// 清空日志
	utils.SharedLogger.Reset()

	this.Data["isOk"] = false

	var installer = mysqlinstallers.NewMySQLInstaller()
	var targetDir = "/usr/local/mysql"
	xzFile, err := installer.Download()
	if err != nil {
		this.Data["err"] = "download failed: " + err.Error()
		this.Success()
		return
	}

	err = installer.InstallFromFile(xzFile, targetDir)
	if err != nil {
		this.Data["err"] = "install from '" + xzFile + "' failed: " + err.Error()
		this.Success()
		return
	}

	this.Data["user"] = "root"
	this.Data["password"] = installer.Password()
	this.Data["dir"] = targetDir
	this.Data["isOk"] = true

	this.Success()
}
