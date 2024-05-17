// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/conds/condutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"io"
	"os"
	"path/filepath"
)

func Generate() error {
	err := generateComponentsJSFile()
	if err != nil {
		return fmt.Errorf("generate 'components.src.js' failed: %w", err)
	}

	return nil
}

// 生成Javascript文件
func generateComponentsJSFile() error {
	var buffer = bytes.NewBuffer([]byte{})

	var webRoot string
	if Tea.IsTesting() {
		webRoot = Tea.Root + "/../web/public/js/components/"
	} else {
		webRoot = Tea.Root + "/web/public/js/components/"
	}
	var f = files.NewFile(webRoot)

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
		buffer.Write([]byte{';', '\n', '\n'})
	}

	// 条件操作符
	requestOperatorsJSON, err := json.Marshal(shared.AllRequestOperators())
	if err != nil {
		logs.Println("ComponentsAction marshal request operators failed: " + err.Error())
	} else {
		buffer.WriteString("window.REQUEST_COND_OPERATORS = ")
		buffer.Write(requestOperatorsJSON)
		buffer.Write([]byte{';', '\n', '\n'})
	}

	// 请求变量
	requestVariablesJSON, err := json.Marshal(shared.DefaultRequestVariables())
	if err != nil {
		logs.Println("ComponentsAction marshal request variables failed: " + err.Error())
	} else {
		buffer.WriteString("window.REQUEST_VARIABLES = ")
		buffer.Write(requestVariablesJSON)
		buffer.Write([]byte{';', '\n', '\n'})
	}

	// 指标
	metricHTTPKeysJSON, err := json.Marshal(serverconfigs.FindAllMetricKeyDefinitions(serverconfigs.MetricItemCategoryHTTP))
	if err != nil {
		logs.Println("ComponentsAction marshal metric http keys failed: " + err.Error())
	} else {
		buffer.WriteString("window.METRIC_HTTP_KEYS = ")
		buffer.Write(metricHTTPKeysJSON)
		buffer.Write([]byte{';', '\n', '\n'})
	}

	// IP地址阈值项目
	ipAddrThresholdItemsJSON, err := json.Marshal(nodeconfigs.FindAllIPAddressThresholdItems())
	if err != nil {
		logs.Println("ComponentsAction marshal ip addr threshold items failed: " + err.Error())
	} else {
		buffer.WriteString("window.IP_ADDR_THRESHOLD_ITEMS = ")
		buffer.Write(ipAddrThresholdItemsJSON)
		buffer.Write([]byte{';', '\n', '\n'})
	}

	// IP地址阈值动作
	ipAddrThresholdActionsJSON, err := json.Marshal(nodeconfigs.FindAllIPAddressThresholdActions())
	if err != nil {
		logs.Println("ComponentsAction marshal ip addr threshold actions failed: " + err.Error())
	} else {
		buffer.WriteString("window.IP_ADDR_THRESHOLD_ACTIONS = ")
		buffer.Write(ipAddrThresholdActionsJSON)
		buffer.Write([]byte{';', '\n', '\n'})
	}

	// WAF checkpoints
	var wafCheckpointsMaps = []maps.Map{}
	for _, checkpoint := range firewallconfigs.AllCheckpoints {
		wafCheckpointsMaps = append(wafCheckpointsMaps, maps.Map{
			"name":        checkpoint.Name,
			"prefix":      checkpoint.Prefix,
			"description": checkpoint.Description,
		})
	}
	wafCheckpointsJSON, err := json.Marshal(wafCheckpointsMaps)
	if err != nil {
		logs.Println("ComponentsAction marshal waf rule checkpoints failed: " + err.Error())
	} else {
		buffer.WriteString("window.WAF_RULE_CHECKPOINTS = ")
		buffer.Write(wafCheckpointsJSON)
		buffer.Write([]byte{';', '\n', '\n'})
	}

	// WAF操作符
	wafOperatorsJSON, err := json.Marshal(firewallconfigs.AllRuleOperators)
	if err != nil {
		logs.Println("ComponentsAction marshal waf rule operators failed: " + err.Error())
	} else {
		buffer.WriteString("window.WAF_RULE_OPERATORS = ")
		buffer.Write(wafOperatorsJSON)
		buffer.Write([]byte{';', '\n', '\n'})
	}

	// WAF验证码类型
	captchaTypesJSON, err := json.Marshal(firewallconfigs.FindAllCaptchaTypes())
	if err != nil {
		logs.Println("ComponentsAction marshal captcha types failed: " + err.Error())
	} else {
		buffer.WriteString("window.WAF_CAPTCHA_TYPES = ")
		buffer.Write(captchaTypesJSON)
		buffer.Write([]byte{';', '\n', '\n'})
	}

	fp, err := os.OpenFile(filepath.Clean(Tea.PublicFile("/js/components.src.js")), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}

	_, err = io.Copy(fp, buffer)
	if err != nil {
		return err
	}

	return nil
}
