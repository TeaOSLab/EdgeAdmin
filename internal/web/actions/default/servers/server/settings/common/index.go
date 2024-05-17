// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package common

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("common")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 只有HTTP服务才支持
	if this.FilterHTTPFamily() {
		return
	}

	this.Data["hasGroupConfig"] = false

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id

	this.Data["commonConfig"] = maps.Map{
		"mergeSlashes": webConfig.MergeSlashes,
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId        int64
	MergeSlashes bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.ServerCommon_LogUpdateCommonSettings, params.WebId)

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebCommon(this.AdminContext(), &pb.UpdateHTTPWebCommonRequest{
		HttpWebId:    params.WebId,
		MergeSlashes: params.MergeSlashes,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
