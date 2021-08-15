// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipbox

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"strings"
	"time"
)

type AddIPAction struct {
	actionutils.ParentAction
}

func (this *AddIPAction) RunPost(params struct {
	ListId int64
	Ip     string
}) {
	var ipType = "ipv4"
	if strings.Contains(params.Ip, ":") {
		ipType = "ipv6"
	}

	_, err := this.RPC().IPItemRPC().CreateIPItem(this.AdminContext(), &pb.CreateIPItemRequest{
		IpListId:   params.ListId,
		IpFrom:     params.Ip,
		IpTo:       "",
		ExpiredAt:  time.Now().Unix() + 86400, // TODO 可以自定义时间
		Reason:     "从IPBox中加入名单",
		Type:       ipType,
		EventLevel: "critical",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
