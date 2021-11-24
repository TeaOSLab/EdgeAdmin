package pages

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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
	this.SecondMenu("pages")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 分组设置
	groupResp, err := this.RPC().ServerGroupRPC().FindEnabledServerGroupConfigInfo(this.AdminContext(), &pb.FindEnabledServerGroupConfigInfoRequest{
		ServerId: params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasGroupConfig"] = groupResp.HasPagesConfig
	this.Data["groupSettingURL"] = "/servers/groups/group/settings/pages?groupId=" + types.String(groupResp.ServerGroupId)

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["pages"] = webConfig.Pages
	this.Data["shutdownConfig"] = webConfig.Shutdown

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId        int64
	PagesJSON    string
	ShutdownJSON string
	Must         *actions.Must
}) {
	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "修改Web %d 的设置", params.WebId)

	// TODO 检查配置

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebPages(this.AdminContext(), &pb.UpdateHTTPWebPagesRequest{
		HttpWebId: params.WebId,
		PagesJSON: []byte(params.PagesJSON),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebShutdown(this.AdminContext(), &pb.UpdateHTTPWebShutdownRequest{
		HttpWebId:    params.WebId,
		ShutdownJSON: []byte(params.ShutdownJSON),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
