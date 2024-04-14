// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package origins

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type UpdateIsOnAction struct {
	actionutils.ParentAction
}

func (this *UpdateIsOnAction) RunPost(params struct {
	OriginId int64
	IsOn     bool
}) {
	defer this.CreateLogInfo(codes.ServerOrigin_LogUpdateOriginIsOn, params.OriginId)

	_, err := this.RPC().OriginRPC().UpdateOriginIsOn(this.AdminContext(), &pb.UpdateOriginIsOnRequest{
		OriginId: params.OriginId,
		IsOn:     params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
