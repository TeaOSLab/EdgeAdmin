// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package addr

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/ip-addrs/ipaddrutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	AddrId int64
}) {
	addr, err := ipaddrutils.InitIPAddr(this.Parent(), params.AddrId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var thresholds = []*nodeconfigs.NodeValueThresholdConfig{}
	if len(addr.ThresholdsJSON) > 0 {
		err = json.Unmarshal(addr.ThresholdsJSON, &thresholds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["supportThresholds"] = true

	this.Data["addr"] = maps.Map{
		"id":          addr.Id,
		"name":        addr.Name,
		"description": addr.Description,
		"ip":          addr.Ip,
		"canAccess":   addr.CanAccess,
		"isOn":        addr.IsOn,
		"isUp":        addr.IsUp,
		"thresholds":  thresholds,
	}

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	AddrId         int64
	IP             string `alias:"ip"`
	Name           string
	CanAccess      bool
	IsOn           bool
	ThresholdsJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("ip", params.IP).
		Require("请输入IP地址")

	ip := net.ParseIP(params.IP)
	if len(ip) == 0 {
		this.Fail("请输入正确的IP")
	}

	_, err := this.RPC().NodeIPAddressRPC().UpdateNodeIPAddress(this.AdminContext(), &pb.UpdateNodeIPAddressRequest{
		NodeIPAddressId: params.AddrId,
		Name:            params.Name,
		Ip:              params.IP,
		CanAccess:       params.CanAccess,
		IsOn:            params.IsOn,
		ThresholdsJSON:  params.ThresholdsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
