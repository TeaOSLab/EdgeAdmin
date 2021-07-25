// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package keys

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	KeyId int64
}) {
	defer this.CreateLogInfo("删除DNS密钥 %d", params.KeyId)

	_, err := this.RPC().NSKeyRPC().DeleteNSKey(this.AdminContext(), &pb.DeleteNSKeyRequest{NsKeyId: params.KeyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
