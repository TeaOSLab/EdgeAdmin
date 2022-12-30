// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package useragent

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("userAgent")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	this.Data["serverId"] = params.ServerId

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id

	var userAgentConfig = webConfig.UserAgent
	if userAgentConfig == nil {
		userAgentConfig = serverconfigs.NewUserAgentConfig()
	}

	this.Data["userAgentConfig"] = userAgentConfig

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId         int64
	UserAgentJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改Web %d User-Agent设置", params.WebId)

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebUserAgent(this.AdminContext(), &pb.UpdateHTTPWebUserAgentRequest{
		HttpWebId:     params.WebId,
		UserAgentJSON: params.UserAgentJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
