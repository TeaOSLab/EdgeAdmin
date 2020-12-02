package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/ttlcache"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"strconv"
)

type ExportDownloadAction struct {
	actionutils.ParentAction
}

func (this *ExportDownloadAction) Init() {
	this.Nav("", "", "")
}

func (this *ExportDownloadAction) RunGet(params struct {
	Key string
}) {
	item := ttlcache.DefaultCache.Read(params.Key)
	if item == nil || item.Value == nil {
		this.WriteString("找不到要导出的内容")
		return
	}

	ttlcache.DefaultCache.Delete(params.Key)

	data, ok := item.Value.([]byte)
	if ok {
		this.AddHeader("Content-Disposition", "attachment; filename=\"WAF.json\";")
		this.AddHeader("Content-Length", strconv.Itoa(len(data)))
		this.Write(data)
	} else {
		this.WriteString("找不到要导出的内容")
		return
	}
}
