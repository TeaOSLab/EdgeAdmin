package ui

import (
	"bytes"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/conds/condutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

type ComponentsAction actions.Action

var componentsData = []byte{}

func (this *ComponentsAction) RunGet(params struct{}) {
	this.AddHeader("Content-Type", "text/javascript; charset=utf-8")

	if !Tea.IsTesting() && len(componentsData) > 0 {
		this.AddHeader("Last-Modified", "Fri, 06 Sep 2019 08:29:50 GMT")
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
		logs.Println("ComponentsAction marshal request cond types failed: " + err.Error())
	} else {
		buffer.WriteString("window.REQUEST_COND_COMPONENTS = ")
		buffer.Write(typesJSON)
		buffer.Write([]byte{'\n', '\n'})
	}

	// 条件操作符
	requestOperatorsJSON, err := json.Marshal(shared.AllRequestOperators())
	if err != nil {
		logs.Println("ComponentsAction marshal request operators failed: " + err.Error())
	} else {
		buffer.WriteString("window.REQUEST_COND_OPERATORS = ")
		buffer.Write(requestOperatorsJSON)
		buffer.Write([]byte{'\n', '\n'})
	}

	// 请求变量
	requestVariablesJSON, err := json.Marshal(shared.DefaultRequestVariables())
	if err != nil {
		logs.Println("ComponentsAction marshal request variables failed: " + err.Error())
	} else {
		buffer.WriteString("window.REQUEST_VARIABLES = ")
		buffer.Write(requestVariablesJSON)
		buffer.Write([]byte{'\n', '\n'})
	}

	componentsData = buffer.Bytes()
	this.Write(componentsData)
}
