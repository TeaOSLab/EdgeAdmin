// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package login

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/index/loginutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/rands"
)

type TicketAction struct {
	actionutils.ParentAction
}

func (this *TicketAction) Init() {
	this.Nav("", "", "")
}

func (this *TicketAction) RunGet(params struct {
	Ticket   string
	Redirect string
	Auth     *helpers.UserShouldAuth
}) {
	this.Data["redirect"] = params.Redirect
	var errorMsg string

	defer func() {
		this.Data["errorMsg"] = errorMsg
		this.Show()
	}()

	if len(params.Ticket) == 0 {
		errorMsg = "invalid ticket: wrong format"
		return
	}

	// TODO 对于错误尝试太多的IP进行处罚

	resp, err := this.RPC().LoginTicketRPC().FindLoginTicketWithValue(this.AdminContext(), &pb.FindLoginTicketWithValueRequest{Value: params.Ticket})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if resp.LoginTicket == nil {
		errorMsg = "invalid ticket: not found"
		return
	}

	if resp.LoginTicket.AdminId <= 0 {
		errorMsg = "invalid ticket: invalid admin id"
		return
	}

	var currentIP = loginutils.RemoteIP(&this.ActionObject)
	if len(resp.LoginTicket.Ip) > 0 && resp.LoginTicket.Ip != currentIP {
		errorMsg = "invalid ticket: wrong client ip"
		return
	}

	var localSid = rands.HexString(32)
	this.Data["localSid"] = localSid
	this.Data["ip"] = currentIP

	params.Auth.StoreAdmin(resp.LoginTicket.AdminId, false, localSid)
}
