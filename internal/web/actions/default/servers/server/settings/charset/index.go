package charset

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/configutils"
	"github.com/iwind/TeaGo/actions"
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
	webConfigResp, err := this.RPC().ServerRPC().FindAndInitServerWebConfig(this.AdminContext(), &pb.FindAndInitServerWebRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	webConfig := &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(webConfigResp.Config, webConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["charset"] = webConfig.Charset

	this.Data["usualCharsets"] = configutils.UsualCharsets
	this.Data["allCharsets"] = configutils.AllCharsets

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId   int64
	Charset string

	Must *actions.Must
}) {
	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebCharset(this.AdminContext(), &pb.UpdateHTTPWebCharsetRequest{
		WebId:   params.WebId,
		Charset: params.Charset,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
