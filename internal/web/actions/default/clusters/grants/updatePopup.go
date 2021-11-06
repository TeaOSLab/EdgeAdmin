package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	GrantId int64
}) {
	this.Data["methods"] = grantutils.AllGrantMethods()

	grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if grantResp.NodeGrant == nil {
		this.WriteString("找不到要操作的对象")
		return
	}

	grant := grantResp.NodeGrant
	this.Data["grant"] = maps.Map{
		"id":          grant.Id,
		"nodeId":      grant.NodeId,
		"method":      grant.Method,
		"name":        grant.Name,
		"username":    grant.Username,
		"password":    grant.Password,
		"description": grant.Description,
		"privateKey":  grant.PrivateKey,
		"passphrase":  grant.Passphrase,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	GrantId     int64
	NodeId      int64
	Name        string
	Method      string
	Username    string
	Password    string
	PrivateKey  string
	Passphrase  string
	Description string

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改SSH认证 %d", params.GrantId)

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

	// 执行修改
	_, err := this.RPC().NodeGrantRPC().UpdateNodeGrant(this.AdminContext(), &pb.UpdateNodeGrantRequest{
		NodeGrantId: params.GrantId,
		Name:        params.Name,
		Method:      params.Method,
		Username:    params.Username,
		Password:    params.Password,
		PrivateKey:  params.PrivateKey,
		Passphrase:  params.Passphrase,
		Description: params.Description,
		NodeId:      params.NodeId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 返回信息
	this.Data["grant"] = maps.Map{
		"id":         params.GrantId,
		"name":       params.Name,
		"method":     params.Method,
		"methodName": grantutils.FindGrantMethodName(params.Method),
	}

	this.Success()
}
