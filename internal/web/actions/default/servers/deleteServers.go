// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

// DeleteServersAction 删除一组网站
type DeleteServersAction struct {
	actionutils.ParentAction
}

func (this *DeleteServersAction) RunPost(params struct {
	ServerIds []int64
}) {
	defer this.CreateLogInfo(codes.Server_LogDeleteServers)

	_, err := this.RPC().ServerRPC().DeleteServers(this.AdminContext(), &pb.DeleteServersRequest{ServerIds: params.ServerIds})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
