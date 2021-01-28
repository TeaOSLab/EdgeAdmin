package ui

import (
	"bytes"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/conds/condutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"strconv"
)

type ComponentsAction actions.Action

var componentsData = []byte{}

func (this *ComponentsAction) RunGet(params struct{}) {
	this.AddHeader("Content-Type", "text/javascript; charset=utf-8")

	if !Tea.IsTesting() && len(componentsData) > 0 {
		this.AddHeader("Content-Length", strconv.Itoa(len(componentsData)))
		this.Write(componentsData)
		return
	}

	var buffer = bytes.NewBuffer([]byte{})

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
		buffer.Write(data)
		buffer.Write([]byte{'\n', '\n'})
	})

	// 条件组件
	typesJSON, err := json.Marshal(condutils.ReadAllAvailableCondTypes())
	if err != nil {
		logs.Println("ComponentsAction: " + err.Error())
	} else {
		buffer.WriteString("window.REQUEST_COND_COMPONENTS = ")
		buffer.Write(typesJSON)
		buffer.Write([]byte{'\n', '\n'})
	}

	componentsData = buffer.Bytes()
	this.AddHeader("Content-Length", strconv.Itoa(len(componentsData)))
	this.Write(componentsData)
}
