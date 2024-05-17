// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package ui

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"os"
	"sync"
)

var tipKeyMap = map[string]bool{}
var tipKeyLocker = sync.Mutex{}
var tipConfigFile = "tip.cache.json"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 从配置文件中加载已关闭的tips
		data, err := os.ReadFile(Tea.ConfigFile(tipConfigFile))
		if err == nil {
			var m = map[string]bool{}
			err = json.Unmarshal(data, &m)
			if err == nil {
				tipKeyLocker.Lock()
				tipKeyMap = m
				tipKeyLocker.Unlock()
			}
		}
	})
}

type ShowTipAction struct {
	actionutils.ParentAction
}

func (this *ShowTipAction) RunPost(params struct {
	Code string
}) {
	tipKeyLocker.Lock()
	_, ok := tipKeyMap[params.Code]
	tipKeyLocker.Unlock()

	this.Data["visible"] = !ok

	this.Success()
}
