package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
)

// 下载指定的文本内容
type DownloadAction struct {
	actionutils.ParentAction
}

func (this *DownloadAction) Init() {
	this.Nav("", "", "")
}

func (this *DownloadAction) RunGet(params struct {
	File string
	Text string

	Auth *helpers.UserMustAuth
}) {
	this.AddHeader("Content-Disposition", "attachment; filename=\""+params.File+"\";")
	this.WriteString(params.Text)
}
