package stat

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
	this.SecondMenu("stat")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
}) {
	_, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "stat")
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
	this.Data["statConfig"] = webConfig.StatRef

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId    int64
	StatJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.ServerStat_LogUpdateStatSettings, params.WebId)

	// TODO 校验配置

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebStat(this.AdminContext(), &pb.UpdateHTTPWebStatRequest{
		HttpWebId: params.WebId,
		StatJSON:  params.StatJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
