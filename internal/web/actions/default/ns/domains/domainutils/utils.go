// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package domainutils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

// InitDomain 初始化域名信息
func InitDomain(parent *actionutils.ParentAction, domainId int64) error {
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	domainResp, err := rpcClient.NSDomainRPC().FindEnabledNSDomain(parent.AdminContext(), &pb.FindEnabledNSDomainRequest{NsDomainId: domainId})
	if err != nil {
		return err
	}
	var domain = domainResp.NsDomain
	if domain == nil {
		return errors.New("InitDomain: can not find domain with id '" + types.String(domainId) + "'")
	}

	// 记录数量
	countRecordsResp, err := rpcClient.NSRecordRPC().CountAllEnabledNSRecords(parent.AdminContext(), &pb.CountAllEnabledNSRecordsRequest{
		NsDomainId: domainId,
	})
	if err != nil {
		return err
	}
	var countRecords = countRecordsResp.Count

	// Key数量
	countKeysResp, err := rpcClient.NSKeyRPC().CountAllEnabledNSKeys(parent.AdminContext(), &pb.CountAllEnabledNSKeysRequest{
		NsDomainId: domainId,
	})
	if err != nil {
		return err
	}
	var countKeys = countKeysResp.Count

	parent.Data["domain"] = maps.Map{
		"id":           domain.Id,
		"name":         domain.Name,
		"countRecords": countRecords,
		"countKeys":    countKeys,
	}

	return nil
}
