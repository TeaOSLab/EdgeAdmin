// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package message

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type SelectedReceiversAction struct {
	actionutils.ParentAction
}

func (this *SelectedReceiversAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectedReceiversAction) RunPost(params struct {
	ClusterId int64
	NodeId    int64
	ServerId  int64
}) {
	receiversResp, err := this.RPC().MessageReceiverRPC().FindAllEnabledMessageReceivers(this.AdminContext(), &pb.FindAllEnabledMessageReceiversRequest{
		NodeClusterId: params.ClusterId,
		NodeId:        params.NodeId,
		ServerId:      params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	receiverMaps := []maps.Map{}
	for _, receiver := range receiversResp.MessageReceivers {
		id := int64(0)
		name := ""
		receiverType := ""
		subName := ""
		if receiver.MessageRecipient != nil {
			id = receiver.MessageRecipient.Id
			name = receiver.MessageRecipient.Admin.Fullname
			subName = receiver.MessageRecipient.MessageMediaInstance.Name
			receiverType = "recipient"
		} else if receiver.MessageRecipientGroup != nil {
			id = receiver.MessageRecipientGroup.Id
			name = receiver.MessageRecipientGroup.Name
			receiverType = "group"
		} else {
			continue
		}
		receiverMaps = append(receiverMaps, maps.Map{
			"id":      id,
			"name":    name,
			"subName": subName,
			"type":    receiverType,
		})
	}
	this.Data["receivers"] = receiverMaps

	this.Success()
}
