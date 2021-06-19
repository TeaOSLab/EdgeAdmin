package access

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("access")
}

func (this *IndexAction) RunGet(params struct {
	LocationId int64
}) {
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithLocationId(this.AdminContext(), params.LocationId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["authConfig"] = webConfig.Auth

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId    int64
	AuthJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改Web %d 的认证设置", params.WebId)

	var authConfig = &serverconfigs.HTTPAuthConfig{}
	err := json.Unmarshal(params.AuthJSON, authConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	err = authConfig.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
	}

	// 保存之前删除多于的配置信息
	for _, ref := range authConfig.PolicyRefs {
		ref.AuthPolicy = nil
	}

	configJSON, err := json.Marshal(authConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebAuth(this.AdminContext(), &pb.UpdateHTTPWebAuthRequest{
		WebId:    params.WebId,
		AuthJSON: configJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
