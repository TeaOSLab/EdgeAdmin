package ui

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/conds/condutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"net/http"
)

type ComponentsAction actions.Action

var componentsData = []byte{}
var componentsDataSum string

func (this *ComponentsAction) RunGet(params struct{}) {
	this.AddHeader("Content-Type", "text/javascript; charset=utf-8")

	// etag
	var requestETag = this.Header("If-None-Match")
	if len(requestETag) > 0 && requestETag == "\""+componentsDataSum+"\"" {
		this.ResponseWriter.WriteHeader(http.StatusNotModified)
		return
	}

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

	// 指标
	metricHTTPKeysJSON, err := json.Marshal(serverconfigs.FindAllMetricKeyDefinitions(serverconfigs.MetricItemCategoryHTTP))
	if err != nil {
		logs.Println("ComponentsAction marshal metric http keys failed: " + err.Error())
	} else {
		buffer.WriteString("window.METRIC_HTTP_KEYS = ")
		buffer.Write(metricHTTPKeysJSON)
		buffer.Write([]byte{'\n', '\n'})
	}

	componentsData = buffer.Bytes()

	// ETag
	var h = md5.New()
	h.Write(buffer.Bytes())
	componentsDataSum = fmt.Sprintf("%x", h.Sum(nil))
	this.AddHeader("ETag", "\""+componentsDataSum+"\"")

	this.Write(componentsData)
}
