// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package common

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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
	defer this.CreateLogInfo("修改服务Web %d 设置的其他设置", params.WebId)

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
