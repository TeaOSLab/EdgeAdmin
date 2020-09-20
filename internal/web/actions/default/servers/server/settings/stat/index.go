package stat

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
	this.SecondMenu("stat")
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
	this.Data["statConfig"] = webConfig.StatRef

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId    int64
	StatJSON []byte

	Must *actions.Must
}) {
	// TODO 校验配置

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPStat(this.AdminContext(), &pb.UpdateHTTPStatRequest{
		WebId:    params.WebId,
		StatJSON: params.StatJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
