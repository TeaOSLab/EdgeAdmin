// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package server

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	keyResp, err := this.RPC().AuthorityKeyRPC().ReadAuthorityKey(this.AdminContext(), &pb.ReadAuthorityKeyRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var keyMap maps.Map = nil
	teaconst.IsPlus = false
	var key = keyResp.AuthorityKey
	if key != nil {
		if len(key.MacAddresses) == 0 {
			key.MacAddresses = []string{}
		}

		isActive := len(key.DayTo) > 0 && key.DayTo >= timeutil.Format("Y-m-d")
		if isActive {
			teaconst.IsPlus = true
		}

		isExpiring := isActive && key.DayTo < timeutil.Format("Y-m-d", time.Now().AddDate(0, 0, 7))

		keyMap = maps.Map{
			"dayFrom":      key.DayFrom,
			"dayTo":        key.DayTo,
			"macAddresses": key.MacAddresses,
			"hostname":     key.Hostname,
			"company":      key.Company,
			"nodes":        key.Nodes,
			"isExpired":    !isActive,
			"isExpiring":   isExpiring,
			"updatedTime":  timeutil.FormatTime("Y-m-d H:i:s", keyResp.AuthorityKey.UpdatedAt),
		}
	}
	this.Data["key"] = keyMap

	// 检查是否有认证节点，如果没有认证节点，则自动生成一个
	countResp, err := this.RPC().AuthorityNodeRPC().CountAllEnabledAuthorityNodes(this.AdminContext(), &pb.CountAllEnabledAuthorityNodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count == 0 {
		_, err = this.RPC().AuthorityNodeRPC().CreateAuthorityNode(this.AdminContext(), &pb.CreateAuthorityNodeRequest{
			Name:        "默认节点",
			Description: "系统自动生成的默认节点",
			IsOn:        true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Show()
}
