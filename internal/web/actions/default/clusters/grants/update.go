package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "grant", "update")
}

func (this *UpdateAction) RunGet(params struct {
	GrantId int64
}) {
	this.Data["methods"] = grantutils.AllGrantMethods()

	grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if grantResp.NodeGrant == nil {
		this.WriteString("can not find the grant")
		return
	}

	// TODO 处理节点专用的认证

	grant := grantResp.NodeGrant
	this.Data["grant"] = maps.Map{
		"id":          grant.Id,
		"name":        grant.Name,
		"method":      grant.Method,
		"methodName":  grantutils.FindGrantMethodName(grant.Method),
		"username":    grant.Username,
		"password":    grant.Password,
		"privateKey":  grant.PrivateKey,
		"description": grant.Description,
		"su":          grant.Su,
	}

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	GrantId     int64
	Name        string
	Method      string
	Username    string
	Password    string
	PrivateKey  string
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

	// TODO 检查grantId是否存在

	_, err := this.RPC().NodeGrantRPC().UpdateNodeGrant(this.AdminContext(), &pb.UpdateNodeGrantRequest{
		NodeGrantId: params.GrantId,
		Name:        params.Name,
		Method:      params.Method,
		Username:    params.Username,
		Password:    params.Password,
		PrivateKey:  params.PrivateKey,
		Description: params.Description,
		NodeId:      0,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
