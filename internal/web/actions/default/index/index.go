package index

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

// 首页（登录页）

// TokenKey 加密用的密钥
var TokenKey = stringutil.Rand(32)

func (this *IndexAction) RunGet(params struct {
	From string

	Auth *helpers.UserShouldAuth
}) {
	// DEMO模式
	this.Data["isDemo"] = teaconst.IsDemoMode

	// 检查系统是否已经配置过
	if !setup.IsConfigured() {
		this.RedirectURL("/setup")
		return
	}

	//// 是否新安装
	if setup.IsNewInstalled() {
		this.RedirectURL("/setup/confirm")
		return
	}

	// 已登录跳转到dashboard
	if params.Auth.IsUser() {
		this.RedirectURL("/dashboard")
		return
	}

	this.Data["isUser"] = false
	this.Data["menu"] = "signIn"

	var timestamp = fmt.Sprintf("%d", time.Now().Unix())
	this.Data["token"] = stringutil.Md5(TokenKey+timestamp) + timestamp
	this.Data["from"] = params.From

	uiConfig, err := configloaders.LoadAdminUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["systemName"] = uiConfig.AdminSystemName
	this.Data["showVersion"] = uiConfig.ShowVersion
	if len(uiConfig.Version) > 0 {
		this.Data["version"] = uiConfig.Version
	} else {
		this.Data["version"] = teaconst.Version
	}
	this.Data["faviconFileId"] = uiConfig.FaviconFileId

	securityConfig, err := configloaders.LoadSecurityConfig()
	if err != nil {
		this.Data["rememberLogin"] = false
	} else {
		this.Data["rememberLogin"] = securityConfig.AllowRememberLogin
	}

	this.Show()
}

// RunPost 提交
func (this *IndexAction) RunPost(params struct {
	Token    string
	Username string
	Password string
	OtpCode  string
	Remember bool

	Must *actions.Must
	Auth *helpers.UserShouldAuth
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Field("password", params.Password).
		Require("请输入密码")

	if params.Password == stringutil.Md5("") {
		this.FailField("password", "请输入密码")
	}

	// 检查token
	if len(params.Token) <= 32 {
		this.Fail("请通过登录页面登录")
	}
	var timestampString = params.Token[32:]
	if stringutil.Md5(TokenKey+timestampString) != params.Token[:32] {
		this.FailField("refresh", "登录页面已过期，请刷新后重试")
	}
	var timestamp = types.Int64(timestampString)
	if timestamp < time.Now().Unix()-1800 {
		this.FailField("refresh", "登录页面已过期，请刷新后重试")
	}

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		this.Fail("服务器出了点小问题：" + err.Error())
		return
	}
	resp, err := rpcClient.AdminRPC().LoginAdmin(rpcClient.Context(0), &pb.LoginAdminRequest{
		Username: params.Username,
		Password: params.Password,
	})

	if err != nil {
		err = dao.SharedLogDAO.CreateAdminLog(rpcClient.Context(0), oplogs.LevelError, this.Request.URL.Path, "登录时发生系统错误："+err.Error(), this.RequestRemoteIP())
		if err != nil {
			utils.PrintError(err)
		}

		actionutils.Fail(this, err)
		return
	}

	if !resp.IsOk {
		err = dao.SharedLogDAO.CreateAdminLog(rpcClient.Context(0), oplogs.LevelWarn, this.Request.URL.Path, "登录失败，用户名："+params.Username, this.RequestRemoteIP())
		if err != nil {
			utils.PrintError(err)
		}

		this.Fail("请输入正确的用户名密码")
		return
	}
	var adminId = resp.AdminId

	// 检查是否支持OTP
	checkOTPResp, err := this.RPC().AdminRPC().CheckAdminOTPWithUsername(this.AdminContext(), &pb.CheckAdminOTPWithUsernameRequest{Username: params.Username})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var requireOTP = checkOTPResp.RequireOTP
	this.Data["requireOTP"] = requireOTP
	if requireOTP {
		this.Data["remember"] = params.Remember

		var sid = this.Session().Sid
		this.Data["sid"] = sid
		_, err = this.RPC().LoginSessionRPC().WriteLoginSessionValue(this.AdminContext(), &pb.WriteLoginSessionValueRequest{
			Sid:   sid + "_otp",
			Key:   "adminId",
			Value: types.String(adminId),
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.Success()
		return
	}

	// 写入SESSION
	params.Auth.StoreAdmin(adminId, params.Remember)

	// 记录日志
	err = dao.SharedLogDAO.CreateAdminLog(rpcClient.Context(adminId), oplogs.LevelInfo, this.Request.URL.Path, "成功登录系统，用户名："+params.Username, this.RequestRemoteIP())
	if err != nil {
		utils.PrintError(err)
	}

	this.Success()
}
