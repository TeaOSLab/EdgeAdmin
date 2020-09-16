package web

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("web")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	server, _, isOk := serverutils.FindServer(&this.ParentAction, params.ServerId)
	if !isOk {
		return
	}
	webId := server.WebId

	webConfig := &serverconfigs.HTTPWebConfig{
		Id:   webId,
		IsOn: true,
	}
	if webId > 0 {
		resp, err := this.RPC().HTTPWebRPC().FindEnabledHTTPWeb(this.AdminContext(), &pb.FindEnabledHTTPWebRequest{WebId: webId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if resp.Web != nil {
			web := resp.Web

			webConfig.Id = webId
			webConfig.IsOn = web.IsOn
			webConfig.Root = web.Root
		}
	}

	this.Data["webConfig"] = webConfig

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId int64
	WebId    int64
	Root     string

	Must *actions.Must
}) {
	if params.WebId <= 0 {
		resp, err := this.RPC().HTTPWebRPC().CreateHTTPWeb(this.AdminContext(), &pb.CreateHTTPWebRequest{
			Root: params.Root,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		webId := resp.WebId
		_, err = this.RPC().ServerRPC().UpdateServerWeb(this.AdminContext(), &pb.UpdateServerWebRequest{
			ServerId: params.ServerId,
			WebId:    webId,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		_, err := this.RPC().HTTPWebRPC().UpdateHTTPWeb(this.AdminContext(), &pb.UpdateHTTPWebRequest{
			WebId: params.WebId,
			Root:  params.Root,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Success()
}
