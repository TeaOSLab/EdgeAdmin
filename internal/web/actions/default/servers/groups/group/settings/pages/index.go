package pages

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/groups/group/servergrouputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("pages")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
}) {
	_, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "pages")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerGroupId(this.AdminContext(), params.GroupId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["enableGlobalPages"] = webConfig.EnableGlobalPages
	this.Data["pages"] = webConfig.Pages
	this.Data["shutdownConfig"] = webConfig.Shutdown

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId             int64
	EnableGlobalPages bool
	PagesJSON         string
	ShutdownJSON      string
	Must              *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.ServerPage_LogUpdatePages, params.WebId)

	// TODO 检查配置

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebGlobalPagesEnabled(this.AdminContext(), &pb.UpdateHTTPWebGlobalPagesEnabledRequest{
		HttpWebId: params.WebId,
		IsEnabled: params.EnableGlobalPages,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebPages(this.AdminContext(), &pb.UpdateHTTPWebPagesRequest{
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
