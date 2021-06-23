// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/golang/protobuf/proto"
	"strconv"
)

type ExportDataAction struct {
	actionutils.ParentAction
}

func (this *ExportDataAction) Init() {
	this.Nav("", "", "")
}

func (this *ExportDataAction) RunGet(params struct {
	ListId int64
}) {
	defer this.CreateLogInfo("导出IP名单 %d", params.ListId)

	resp := &pb.ListIPItemsWithListIdResponse{}
	var offset int64 = 0
	var size int64 = 1000
	for {
		itemsResp, err := this.RPC().IPItemRPC().ListIPItemsWithListId(this.AdminContext(), &pb.ListIPItemsWithListIdRequest{
			IpListId: params.ListId,
			Offset:   offset,
			Size:     size,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(itemsResp.IpItems) == 0 {
			break
		}
		for _, item := range itemsResp.IpItems {
			resp.IpItems = append(resp.IpItems, item)
		}
		offset += size
	}

	data, err := proto.Marshal(resp)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.AddHeader("Content-Disposition", "attachment; filename=\"ip-list-"+numberutils.FormatInt64(params.ListId)+".data\";")
	this.AddHeader("Content-Length", strconv.Itoa(len(data)))
	this.Write(data)
}
