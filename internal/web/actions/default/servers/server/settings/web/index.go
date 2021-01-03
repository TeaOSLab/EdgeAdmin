package web

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
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
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["rootConfig"] = webConfig.Root

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId int64
	WebId    int64
	RootJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo("修改Web %d 的首页文件名", params.WebId)

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWeb(this.AdminContext(), &pb.UpdateHTTPWebRequest{
		WebId:    params.WebId,
		RootJSON: params.RootJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
