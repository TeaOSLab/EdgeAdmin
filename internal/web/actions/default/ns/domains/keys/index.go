// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package keys

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/ns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "key")
}

func (this *IndexAction) RunGet(params struct {
	DomainId int64
}) {
	// 初始化域名信息
	err := domainutils.InitDomain(this.Parent(), params.DomainId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 数量
	countResp, err := this.RPC().NSKeyRPC().CountAllEnabledNSKeys(this.AdminContext(), &pb.CountAllEnabledNSKeysRequest{
		NsDomainId: params.DomainId,
		NsZoneId:   0,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var page = this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	// 列表
	keysResp, err := this.RPC().NSKeyRPC().ListEnabledNSKeys(this.AdminContext(), &pb.ListEnabledNSKeysRequest{
		NsDomainId: params.DomainId,
		NsZoneId:   0,
		Offset:     page.Offset,
		Size:       page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var keyMaps = []maps.Map{}
	for _, key := range keysResp.NsKeys {
		keyMaps = append(keyMaps, maps.Map{
			"id":             key.Id,
			"name":           key.Name,
			"secret":         key.Secret,
			"secretTypeName": dnsconfigs.FindKeySecretTypeName(key.SecretType),
			"algoName":       dnsconfigs.FindKeyAlgorithmTypeName(key.Algo),
			"isOn":           key.IsOn,
		})
	}
	this.Data["keys"] = keyMaps

	this.Show()
}
