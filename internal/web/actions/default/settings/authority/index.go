// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package server

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	keyResp, err := this.RPC().AuthorityKeyRPC().ReadAuthorityKey(this.AdminContext(), &pb.ReadAuthorityKeyRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var keyMap maps.Map = nil
	teaconst.IsPlus = false
	if keyResp.AuthorityKey != nil {
		if len(keyResp.AuthorityKey.MacAddresses) == 0 {
			keyResp.AuthorityKey.MacAddresses = []string{}
		}

		isActive := len(keyResp.AuthorityKey.DayTo) > 0 && keyResp.AuthorityKey.DayTo >= timeutil.Format("Y-m-d")
		if isActive {
			teaconst.IsPlus = true
		}

		keyMap = maps.Map{
			"dayFrom":      keyResp.AuthorityKey.DayFrom,
			"dayTo":        keyResp.AuthorityKey.DayTo,
			"macAddresses": keyResp.AuthorityKey.MacAddresses,
			"hostname":     keyResp.AuthorityKey.Hostname,
			"isExpired":    !isActive,
			"updatedTime":  timeutil.FormatTime("Y-m-d H:i:s", keyResp.AuthorityKey.UpdatedAt),
		}
	}
	this.Data["key"] = keyMap

	this.Show()
}
