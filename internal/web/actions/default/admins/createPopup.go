package admins

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/xlzd/gotp"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Data["modules"] = configloaders.AllModuleMaps(this.LangCode())
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Fullname    string
	Username    string
	Pass1       string
	Pass2       string
	ModuleCodes []string
	IsSuper     bool
	CanLogin    bool

	// OTP
	OtpOn bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("fullname", params.Fullname).
		Require("请输入系统用户全名")

	params.Must.
		Field("username", params.Username).
		Require("请输入登录用户名").
		Match(`^[a-zA-Z0-9_]+$`, "用户名中只能包含英文、数字或下划线")

	existsResp, err := this.RPC().AdminRPC().CheckAdminUsername(this.AdminContext(), &pb.CheckAdminUsernameRequest{
		AdminId:  0,
		Username: params.Username,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if existsResp.Exists {
		this.FailField("username", "此用户名已经被别的系统用户使用，请换一个")
	}

	params.Must.
		Field("pass1", params.Pass1).
		Require("请输入登录密码").
		Field("pass2", params.Pass2).
		Require("请输入确认登录密码")
	if params.Pass1 != params.Pass2 {
		this.FailField("pass2", "两次输入的密码不一致")
	}

	modules := []*systemconfigs.AdminModule{}
	for _, code := range params.ModuleCodes {
		modules = append(modules, &systemconfigs.AdminModule{
			Code:     code,
			AllowAll: true,
			Actions:  nil, // TODO 后期再开放细粒度控制
		})
	}
	modulesJSON, err := json.Marshal(modules)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	createResp, err := this.RPC().AdminRPC().CreateAdmin(this.AdminContext(), &pb.CreateAdminRequest{
		Username:    params.Username,
		Password:    params.Pass1,
		Fullname:    params.Fullname,
		ModulesJSON: modulesJSON,
		IsSuper:     params.IsSuper,
		CanLogin:    params.CanLogin,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// OTP
	if params.OtpOn {
		_, err = this.RPC().LoginRPC().UpdateLogin(this.AdminContext(), &pb.UpdateLoginRequest{Login: &pb.Login{
			Id:   0,
			Type: "otp",
			ParamsJSON: maps.Map{
				"secret": gotp.RandomSecret(16), // TODO 改成可以设置secret长度
			}.AsJSON(),
			IsOn:    true,
			AdminId: createResp.AdminId,
			UserId:  0,
		}})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	defer this.CreateLogInfo(codes.Admin_LogCreateAdmin, createResp.AdminId)

	// 通知更改
	err = configloaders.NotifyAdminModuleMappingChange()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
