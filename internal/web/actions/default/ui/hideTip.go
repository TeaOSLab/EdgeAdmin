// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ui

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
)

type HideTipAction struct {
	actionutils.ParentAction
}

func (this *HideTipAction) RunPost(params struct {
	Code string
}) {
	tipKeyLocker.Lock()
	tipKeyMap[params.Code] = true
	tipKeyLocker.Unlock()

	// 保存到文件
	tipJSON, err := json.Marshal(tipKeyMap)
	if err == nil {
		_ = ioutil.WriteFile(Tea.ConfigFile(tipConfigFile), tipJSON, 0666)
	}

	this.Success()
}
