// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package server

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type ActivateAction struct {
	actionutils.ParentAction
}

func (this *ActivateAction) Init() {
	this.Nav("", "", "activate")
}

func (this *ActivateAction) RunGet(params struct{}) {
	this.Show()
}

func (this *ActivateAction) RunPost(params struct {
	Key string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	if len(params.Key) == 0 {
		this.FailField("key", "请输入激活码")
	}

	resp, err := this.RPC().AuthorityKeyRPC().ValidateAuthorityKey(this.AdminContext(), &pb.ValidateAuthorityKeyRequest{Key: params.Key})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if resp.IsOk {
		_, err := this.RPC().AuthorityKeyRPC().UpdateAuthorityKey(this.AdminContext(), &pb.UpdateAuthorityKeyRequest{
			Value: params.Key,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		this.Success()
	} else {
		this.FailField("key", "无法激活："+resp.Error)
	}
}
