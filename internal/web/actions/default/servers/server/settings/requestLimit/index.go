// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package requestlimit

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("requestLimit")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 只有HTTP服务才支持
	if this.FilterHTTPFamily() {
		return
	}

	// 服务分组设置
	groupResp, err := this.RPC().ServerGroupRPC().FindEnabledServerGroupConfigInfo(this.AdminContext(), &pb.FindEnabledServerGroupConfigInfoRequest{
		ServerId: params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasGroupConfig"] = groupResp.HasRequestLimitConfig
	this.Data["groupSettingURL"] = "/servers/groups/group/settings/requestLimit?groupId=" + types.String(groupResp.ServerGroupId)

	this.Data["serverId"] = params.ServerId

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
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
