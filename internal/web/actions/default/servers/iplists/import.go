// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/golang/protobuf/proto"
	"github.com/iwind/TeaGo/actions"
)

type ImportAction struct {
	actionutils.ParentAction
}

func (this *ImportAction) Init() {
	this.Nav("", "", "import")
}

func (this *ImportAction) RunGet(params struct {
	ListId int64
}) {
	err := InitIPList(this.Parent(), params.ListId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}

func (this *ImportAction) RunPost(params struct {
	ListId int64
	File   *actions.File

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("导入IP名单 %d", params.ListId)

	existsResp, err := this.RPC().IPListRPC().ExistsEnabledIPList(this.AdminContext(), &pb.ExistsEnabledIPListRequest{IpListId: params.ListId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if !existsResp.Exists {
		this.Fail("IP名单不存在")
	}

	if params.File == nil {
		this.Fail("请选择要导入的IP文件")
	}

	data, err := params.File.Read()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	resp := &pb.ListIPItemsWithListIdResponse{}
	err = proto.Unmarshal(data, resp)
	if err != nil {
		this.Fail("导入失败，文件格式错误：" + err.Error())
	}

	var count = 0
	var countIgnore = 0
	for _, item := range resp.IpItems {
		_, err = this.RPC().IPItemRPC().CreateIPItem(this.AdminContext(), &pb.CreateIPItemRequest{
			IpListId:   params.ListId,
			IpFrom:     item.IpFrom,
			IpTo:       item.IpTo,
			ExpiredAt:  item.ExpiredAt,
			Reason:     item.Reason,
			Type:       item.Type,
			EventLevel: item.EventLevel,
		})
		if err != nil {
			this.Fail("导入过程中出错：" + err.Error())
		}
		count++
	}

	this.Data["count"] = count
	this.Data["countIgnore"] = countIgnore

	this.Success()
}
