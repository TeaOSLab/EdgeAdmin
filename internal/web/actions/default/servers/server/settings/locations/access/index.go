package access

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
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

	// 移除不存在的鉴权方法
	var allTypes = []string{}
	for _, def := range serverconfigs.FindAllHTTPAuthTypes(teaconst.Role) {
		allTypes = append(allTypes, def.Code)
	}

	if webConfig.Auth != nil {
		var refs = webConfig.Auth.PolicyRefs
		var realRefs = []*serverconfigs.HTTPAuthPolicyRef{}
		for _, ref := range refs {
			if ref.AuthPolicy == nil {
				continue
			}
			if !lists.ContainsString(allTypes, ref.AuthPolicy.Type) {
				continue
			}
			realRefs = append(realRefs, ref)
		}
		webConfig.Auth.PolicyRefs = realRefs
	}

	this.Data["authConfig"] = webConfig.Auth

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId    int64
	AuthJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.ServerAuth_LogUpdateHTTPAuthSettings, params.WebId)

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
		HttpWebId: params.WebId,
		AuthJSON:  configJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
