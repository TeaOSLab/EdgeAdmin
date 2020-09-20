package pages

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
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
	ServerId int64
}) {
	webResp, err := this.RPC().ServerRPC().FindAndInitServerWebConfig(this.AdminContext(), &pb.FindAndInitServerWebRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	webConfig := &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(webResp.Config, webConfig)
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
	// TODO 检查配置

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebPages(this.AdminContext(), &pb.UpdateHTTPWebPagesRequest{
		WebId:     params.WebId,
		PagesJSON: []byte(params.PagesJSON),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebShutdown(this.AdminContext(), &pb.UpdateHTTPWebShutdownRequest{
		WebId:        params.WebId,
		ShutdownJSON: []byte(params.ShutdownJSON),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
