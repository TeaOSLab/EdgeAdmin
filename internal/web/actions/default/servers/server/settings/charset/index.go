package charset

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
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
	this.SecondMenu("charset")
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
	this.Data["hasGroupConfig"] = groupResp.HasCharsetConfig
	this.Data["groupSettingURL"] = "/servers/groups/group/settings/charset?groupId=" + types.String(groupResp.ServerGroupId)

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["charsetConfig"] = webConfig.Charset

	this.Data["usualCharsets"] = configutils.UsualCharsets
	this.Data["allCharsets"] = configutils.AllCharsets

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId       int64
	CharsetJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.ServerCharset_LogUpdateCharsetSetting, params.WebId)

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebCharset(this.AdminContext(), &pb.UpdateHTTPWebCharsetRequest{
		HttpWebId:   params.WebId,
		CharsetJSON: params.CharsetJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
