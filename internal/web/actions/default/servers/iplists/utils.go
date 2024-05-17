// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/ipconfigs"
	"github.com/iwind/TeaGo/maps"
)

func InitIPList(action *actionutils.ParentAction, listId int64) error {
	client, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	listResp, err := client.IPListRPC().FindEnabledIPList(action.AdminContext(), &pb.FindEnabledIPListRequest{IpListId: listId})
	if err != nil {
		return err
	}
	list := listResp.IpList
	if list == nil {
		return errors.New("not found")
	}

	var typeName = ""
	switch list.Type {
	case ipconfigs.IPListTypeBlack:
		typeName = "黑名单"
	case ipconfigs.IPListTypeWhite:
		typeName = "白名单"
	case ipconfigs.IPListTypeGrey:
		typeName = "灰名单"
	}

	// IP数量
	countItemsResp, err := client.IPItemRPC().CountIPItemsWithListId(action.AdminContext(), &pb.CountIPItemsWithListIdRequest{IpListId: listId})
	if err != nil {
		return err
	}
	countItems := countItemsResp.Count

	action.Data["list"] = maps.Map{
		"id":          list.Id,
		"name":        list.Name,
		"code":        list.Code,
		"type":        list.Type,
		"typeName":    typeName,
		"description": list.Description,
		"isOn":        list.IsOn,
		"countItems":  countItems,
		"isGlobal":    list.IsGlobal,
	}
	return nil
}
