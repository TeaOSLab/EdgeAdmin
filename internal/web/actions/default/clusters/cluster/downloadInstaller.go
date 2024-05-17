// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/types"
	"io"
	"net/http"
	"os"
	"regexp"
)

type DownloadInstallerAction struct {
	actionutils.ParentAction
}

func (this *DownloadInstallerAction) Init() {
	this.Nav("", "", "")
}

func (this *DownloadInstallerAction) RunGet(params struct {
	Name string
}) {
	if len(params.Name) == 0 {
		this.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.WriteString("file not found")
		return
	}

	// 检查文件名
	// 以防止路径穿越等风险
	if !regexp.MustCompile(`^[a-zA-Z0-9.-]+$`).MatchString(params.Name) {
		this.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.WriteString("file not found")
		return
	}

	var zipFile = Tea.Root + "/edge-api/deploy/" + params.Name
	fp, err := os.OpenFile(zipFile, os.O_RDWR, 0444)
	if err != nil {
		if os.IsNotExist(err) {
			this.ResponseWriter.WriteHeader(http.StatusNotFound)
			this.WriteString("file not found")
			return
		}

		this.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		this.WriteString("file can not be opened")

		return
	}

	defer func() {
		_ = fp.Close()
	}()

	stat, err := fp.Stat()
	if err != nil {
		this.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		this.WriteString("file can not be opened")
		return
	}

	this.AddHeader("Content-Disposition", "attachment; filename=\""+params.Name+"\";")
	this.AddHeader("Content-Type", "application/zip")
	this.AddHeader("Content-Length", types.String(stat.Size()))
	_, _ = io.Copy(this.ResponseWriter, fp)
}
