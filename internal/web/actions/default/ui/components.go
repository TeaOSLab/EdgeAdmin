package ui

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/conds/condutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

type ComponentsAction actions.Action

func (this *ComponentsAction) RunGet(params struct{}) {
	this.AddHeader("Content-Type", "text/javascript; charset=utf-8")

	var webRoot string
	if Tea.IsTesting() {
		webRoot = Tea.Root + "/../web/public/js/components/"
	} else {
		webRoot = Tea.Root + "/web/public/js/components/"
	}
	f := files.NewFile(webRoot)

	f.Range(func(file *files.File) {
		if !file.IsFile() {
			return
		}
		if file.Ext() != ".js" {
			return
		}
		data, err := file.ReadAll()
		if err != nil {
			logs.Error(err)
			return
		}
		this.Write(data)
		this.Write([]byte{'\n', '\n'})
	})

	// 条件组件
	typesJSON, err := json.Marshal(condutils.ReadAllAvailableCondTypes())
	if err != nil {
		logs.Println("ComponentsAction: " + err.Error())
	} else {
		this.WriteString("window.REQUEST_COND_COMPONENTS = ")
		this.Write(typesJSON)
		this.Write([]byte{'\n', '\n'})
	}
}
