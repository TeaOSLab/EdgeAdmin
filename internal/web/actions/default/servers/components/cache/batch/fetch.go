// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/cache/cacheutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

type FetchAction struct {
	actionutils.ParentAction
}

func (this *FetchAction) Init() {
	this.Nav("", "", "fetch")
}

func (this *FetchAction) RunGet(params struct{}) {
	// 初始化菜单数据
	err := InitMenu(this.Parent())
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}

func (this *FetchAction) RunPost(params struct {
	Keys string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.HTTPCacheTask_LogCreateHTTPCacheTaskFetch)

	if len(params.Keys) == 0 {
		this.Fail("请输入要预热的Key列表")
	}

	// 检查Key
	var realKeys = []string{}
	for _, key := range strings.Split(params.Keys, "\n") {
		key = strings.TrimSpace(key)
		if len(key) == 0 {
			continue
		}
		if lists.ContainsString(realKeys, key) {
			continue
		}
		realKeys = append(realKeys, key)
	}

	if len(realKeys) == 0 {
		this.Fail("请输入要预热的Key列表")
	}

	// 校验Key
	validateResp, err := this.RPC().HTTPCacheTaskKeyRPC().ValidateHTTPCacheTaskKeys(this.AdminContext(), &pb.ValidateHTTPCacheTaskKeysRequest{Keys: realKeys})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var failKeyMaps = []maps.Map{}
	if len(validateResp.FailKeys) > 0 {
		for _, key := range validateResp.FailKeys {
			failKeyMaps = append(failKeyMaps, maps.Map{
				"key":    key.Key,
				"reason": cacheutils.KeyFailReason(key.ReasonCode),
			})
		}
	}
	this.Data["failKeys"] = failKeyMaps
	if len(failKeyMaps) > 0 {
		this.Fail("有" + types.String(len(failKeyMaps)) + "个Key无法完成操作，请删除后重试")
	}

	// 提交任务
	_, err = this.RPC().HTTPCacheTaskRPC().CreateHTTPCacheTask(this.AdminContext(), &pb.CreateHTTPCacheTaskRequest{
		Type:    "fetch",
		KeyType: "key",
		Keys:    realKeys,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
