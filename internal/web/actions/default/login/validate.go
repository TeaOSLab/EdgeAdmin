// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package login

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/index/loginutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
	"net"
)

type ValidateAction struct {
	actionutils.ParentAction
}

func (this *ValidateAction) Init() {
	this.Nav("", "", "")
}

func (this *ValidateAction) RunGet(params struct {
	From string
}) {
	this.Data["from"] = params.From

	this.Show()
}

func (this *ValidateAction) RunPost(params struct {
	Must *actions.Must

	LocalSid string
	Ip       string
}) {
	var isOk bool

	defer func() {
		this.Data["isOk"] = isOk

		if !isOk {
			loginutils.UnsetCookie(&this.ActionObject)
			this.Session().Delete()
		}

		this.Success()
	}()

	if len(params.LocalSid) == 0 || len(params.LocalSid) != 32 {
		return
	}
	if len(params.Ip) == 0 {
		return
	}

	if net.ParseIP(params.Ip) == nil {
		return
	}

	if params.LocalSid == this.Session().GetString("@localSid") {
		isOk = true

		// renew ip and local sid
		var newIP = loginutils.RemoteIP(&this.ActionObject)
		var newLocalSid = rands.HexString(32)

		this.Session().Write("@ip", newIP)
		this.Session().Write("@localSid", newLocalSid)

		this.Data["ip"] = newIP
		this.Data["localSid"] = newLocalSid

		return
	}
}
