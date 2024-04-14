package accessLog

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("accessLog")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 只有HTTP服务才支持
	if this.FilterHTTPFamily() {
		return
	}

	// 网站分组设置
	groupResp, err := this.RPC().ServerGroupRPC().FindEnabledServerGroupConfigInfo(this.AdminContext(), &pb.FindEnabledServerGroupConfigInfoRequest{
		ServerId: params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasGroupConfig"] = groupResp.HasAccessLogConfig
	this.Data["groupSettingURL"] = "/servers/groups/group/settings/accessLog?groupId=" + types.String(groupResp.ServerGroupId)

	// 获取配置
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["webId"] = webConfig.Id
	this.Data["accessLogConfig"] = webConfig.AccessLogRef

	// 通用变量
	this.Data["fields"] = serverconfigs.HTTPAccessLogFields
	this.Data["defaultFieldCodes"] = serverconfigs.HTTPAccessLogDefaultFieldsCodes

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId         int64
	AccessLogJSON []byte

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.ServerAccessLog_LogUpdateAccessLogSetting, params.WebId)

	// TODO 检查参数

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebAccessLog(this.AdminContext(), &pb.UpdateHTTPWebAccessLogRequest{
		HttpWebId:     params.WebId,
		AccessLogJSON: params.AccessLogJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
