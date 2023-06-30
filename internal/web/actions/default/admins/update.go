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

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	AdminId int64
}) {
	adminResp, err := this.RPC().AdminRPC().FindEnabledAdmin(this.AdminContext(), &pb.FindEnabledAdminRequest{AdminId: params.AdminId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var admin = adminResp.Admin
	if admin == nil {
		this.NotFound("admin", params.AdminId)
		return
	}

	// OTP认证
	var otpLoginIsOn = false
	if admin.OtpLogin != nil {
		otpLoginIsOn = admin.OtpLogin.IsOn
	}

	// AccessKey数量
	countAccessKeyResp, err := this.RPC().UserAccessKeyRPC().CountAllEnabledUserAccessKeys(this.AdminContext(), &pb.CountAllEnabledUserAccessKeysRequest{AdminId: params.AdminId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var countAccessKeys = countAccessKeyResp.Count

	this.Data["admin"] = maps.Map{
		"id":              admin.Id,
		"fullname":        admin.Fullname,
		"username":        admin.Username,
		"isOn":            admin.IsOn,
		"isSuper":         admin.IsSuper,
		"canLogin":        admin.CanLogin,
		"otpLoginIsOn":    otpLoginIsOn,
		"countAccessKeys": countAccessKeys,
	}

	// 权限
	var moduleMaps = configloaders.AllModuleMaps(this.LangCode())
	for _, m := range moduleMaps {
		code := m.GetString("code")
		isChecked := false
		for _, module := range admin.Modules {
			if module.Code == code {
				isChecked = true
				break
			}
		}
		m["isChecked"] = isChecked
	}
	this.Data["modules"] = moduleMaps

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	AdminId int64

	Fullname    string
	Username    string
	Pass1       string
	Pass2       string
	ModuleCodes []string
	IsOn        bool
	IsSuper     bool
	CanLogin    bool

	// OTP
	OtpOn bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.Admin_LogUpdateAdmin, params.AdminId)

	params.Must.
		Field("fullname", params.Fullname).
		Require("请输入系统用户全名")

	params.Must.
		Field("username", params.Username).
		Require("请输入登录用户名").
		Match(`^[a-zA-Z0-9_]+$`, "用户名中只能包含英文、数字或下划线")

	existsResp, err := this.RPC().AdminRPC().CheckAdminUsername(this.AdminContext(), &pb.CheckAdminUsernameRequest{
		AdminId:  params.AdminId,
		Username: params.Username,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if existsResp.Exists {
		this.FailField("username", "此用户名已经被别的系统用户使用，请换一个")
	}

	if len(params.Pass1) > 0 {
		params.Must.
			Field("pass1", params.Pass1).
			Require("请输入登录密码").
			Field("pass2", params.Pass2).
			Require("请输入确认登录密码")
		if params.Pass1 != params.Pass2 {
			this.FailField("pass2", "两次输入的密码不一致")
		}
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

	_, err = this.RPC().AdminRPC().UpdateAdmin(this.AdminContext(), &pb.UpdateAdminRequest{
		AdminId:     params.AdminId,
		Username:    params.Username,
		Password:    params.Pass1,
		Fullname:    params.Fullname,
		ModulesJSON: modulesJSON,
		IsSuper:     params.IsSuper,
		IsOn:        params.IsOn,
		CanLogin:    params.CanLogin,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 修改OTP
	otpLoginResp, err := this.RPC().LoginRPC().FindEnabledLogin(this.AdminContext(), &pb.FindEnabledLoginRequest{
		AdminId: params.AdminId,
		Type:    "otp",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	{
		otpLogin := otpLoginResp.Login
		if params.OtpOn {
			if otpLogin == nil {
				otpLogin = &pb.Login{
					Id:   0,
					Type: "otp",
					ParamsJSON: maps.Map{
						"secret": gotp.RandomSecret(16), // TODO 改成可以设置secret长度
					}.AsJSON(),
					IsOn:    true,
					AdminId: params.AdminId,
					UserId:  0,
				}
			} else {
				// 如果已经有了，就覆盖，这样可以保留既有的参数
				otpLogin.IsOn = true
			}

			_, err = this.RPC().LoginRPC().UpdateLogin(this.AdminContext(), &pb.UpdateLoginRequest{Login: otpLogin})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		} else {
			_, err = this.RPC().LoginRPC().UpdateLogin(this.AdminContext(), &pb.UpdateLoginRequest{Login: &pb.Login{
				Type:    "otp",
				IsOn:    false,
				AdminId: params.AdminId,
			}})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		// 通知更改
		err = configloaders.NotifyAdminModuleMappingChange()
		if err != nil {
			this.ErrorPage(err)
			return
		}

		this.Success()
	}
}
