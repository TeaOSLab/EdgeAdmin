// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipaddrutils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

func InitIPAddr(parent *actionutils.ParentAction, addrId int64) (*pb.NodeIPAddress, error) {
	addrResp, err := parent.RPC().NodeIPAddressRPC().FindEnabledNodeIPAddress(parent.AdminContext(), &pb.FindEnabledNodeIPAddressRequest{NodeIPAddressId: addrId})
	if err != nil {
		return nil, err
	}
	var addr = addrResp.NodeIPAddress
	if addr == nil {
		return nil, errors.New("nodeIPAddress with id '" + types.String(addrId) + "' not found")
	}

	parent.Data["addr"] = maps.Map{
		"id":   addr.Id,
		"name": addr.Name,
		"ip":   addr.Ip,
	}

	return addr, nil
}
