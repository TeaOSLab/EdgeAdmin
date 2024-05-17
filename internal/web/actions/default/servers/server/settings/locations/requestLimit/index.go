// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package requestlimit

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {

}

func (this *IndexAction) RunGet(params struct {
	ServerId   int64
	LocationId int64
}) {
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithLocationId(this.AdminContext(), params.LocationId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["requestLimitConfig"] = webConfig.RequestLimit

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId            int64
	RequestLimitJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.ServerRequestLimit_LogUpdateRequestLimitSettings, params.WebId)

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebRequestLimit(this.AdminContext(), &pb.UpdateHTTPWebRequestLimitRequest{
		HttpWebId:        params.WebId,
		RequestLimitJSON: params.RequestLimitJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
