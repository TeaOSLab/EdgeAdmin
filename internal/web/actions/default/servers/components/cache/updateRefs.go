// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type UpdateRefsAction struct {
	actionutils.ParentAction
}

func (this *UpdateRefsAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateRefsAction) RunPost(params struct {
	CachePolicyId int64
	RefsJSON      []byte
}) {
	// 修改缓存条件
	if params.CachePolicyId > 0 && len(params.RefsJSON) > 0 {
		_, err := this.RPC().HTTPCachePolicyRPC().UpdateHTTPCachePolicyRefs(this.AdminContext(), &pb.UpdateHTTPCachePolicyRefsRequest{
			HttpCachePolicyId: params.CachePolicyId,
			RefsJSON:          params.RefsJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Success()
}
