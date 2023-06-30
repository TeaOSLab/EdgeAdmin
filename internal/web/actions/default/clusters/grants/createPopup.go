package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"golang.org/x/crypto/ssh"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Data["methods"] = grantutils.AllGrantMethods(this.LangCode())

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
	Su          bool

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

		// 验证私钥
		var err error
		if len(params.Passphrase) > 0 {
			_, err = ssh.ParsePrivateKeyWithPassphrase([]byte(params.PrivateKey), []byte(params.Passphrase))
		} else {
			_, err = ssh.ParsePrivateKey([]byte(params.PrivateKey))
		}
		if err != nil {
			this.Fail("私钥验证失败，请检查格式：" + err.Error())
			return
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
		Su:          params.Su,
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
		"methodName": grantutils.FindGrantMethodName(params.Method, this.LangCode()),
		"username":   params.Username,
	}

	// 创建日志
	defer this.CreateLogInfo(codes.NodeGrant_LogCreateSSHGrant, createResp.NodeGrantId)

	this.Success()
}
