// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipbox

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/accesslogs/policyutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

type PolicyAction struct {
	actionutils.ParentAction
}

func (this *PolicyAction) Init() {
	this.Nav("", "", "policy")
}

func (this *PolicyAction) RunGet(params struct {
	PolicyId int64
}) {
	err := policyutils.InitPolicy(this.Parent(), params.PolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var policyMap = this.Data.GetMap("policy")
	if policyMap.GetString("type") == serverconfigs.AccessLogStorageTypeSyslog {
		this.Data["syslogPriorityName"] = serverconfigs.FindAccessLogSyslogStoragePriorityName(policyMap.GetMap("options").GetInt("priority"))
	} else {
		this.Data["syslogPriorityName"] = ""
	}

	this.Show()
}
