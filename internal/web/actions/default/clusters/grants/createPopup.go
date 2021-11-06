package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Data["methods"] = grantutils.AllGrantMethods()

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name        string
	Method      string
	Username    string
	Password    string
	PrivateKey  string
	Passphrase  string
	Description string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入名称")

	switch params.Method {
	case "user":
		if len(params.Username) == 0 {
			this.FailField("username", "请输入SSH登录用户名")
		}
	case "privateKey":
		if len(params.Username) == 0 {
			this.FailField("username", "请输入SSH登录用户名")
		}
		if len(params.PrivateKey) == 0 {
			this.FailField("privateKey", "请输入RSA私钥")
		}
	default:
		this.Fail("请选择正确的认证方式")
	}

	createResp, err := this.RPC().NodeGrantRPC().CreateNodeGrant(this.AdminContext(), &pb.CreateNodeGrantRequest{
		Name:        params.Name,
		Method:      params.Method,
		Username:    params.Username,
		Password:    params.Password,
		PrivateKey:  params.PrivateKey,
		Passphrase:  params.Passphrase,
		Description: params.Description,
		NodeId:      0,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["grant"] = maps.Map{
		"id":         createResp.NodeGrantId,
		"name":       params.Name,
		"method":     params.Method,
		"methodName": grantutils.FindGrantMethodName(params.Method),
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "创建SSH认证 %d", createResp.NodeGrantId)

	this.Success()
}
