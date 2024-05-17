// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package index

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/index/loginutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"github.com/xlzd/gotp"
	"time"
)

type OtpAction struct {
	actionutils.ParentAction
}

func (this *OtpAction) Init() {
	this.Nav("", "", "")
}

func (this *OtpAction) RunGet(params struct {
	From     string
	Sid      string
	Remember bool
}) {
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

	this.Data["isUser"] = false
	this.Data["menu"] = "signIn"

	var timestamp = fmt.Sprintf("%d", time.Now().Unix())
	this.Data["token"] = stringutil.Md5(TokenKey+timestamp) + timestamp
	this.Data["from"] = params.From
	this.Data["sid"] = params.Sid

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
	this.Data["remember"] = params.Remember

	this.Show()
}

func (this *OtpAction) RunPost(params struct {
	Sid      string
	OtpCode  string
	Remember bool

	Must *actions.Must
	Auth *helpers.UserShouldAuth
}) {
	if len(params.OtpCode) == 0 {
		this.FailField("otpCode", "请输入正确的OTP动态密码")
		return
	}

	var sid = params.Sid
	if len(sid) == 0 || len(sid) > 64 {
		this.Fail("参数错误，请重新登录（001）")
		return
	}
	sid += "_otp"

	// 获取SESSION
	sessionResp, err := this.RPC().LoginSessionRPC().FindLoginSession(this.AdminContext(), &pb.FindLoginSessionRequest{Sid: sid})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var session = sessionResp.LoginSession
	if session == nil || session.AdminId <= 0 {
		this.Fail("参数错误，请重新登录（002）")
		return
	}
	var adminId = session.AdminId

	// 检查OTP
	otpLoginResp, err := this.RPC().LoginRPC().FindEnabledLogin(this.AdminContext(), &pb.FindEnabledLoginRequest{
		AdminId: adminId,
		Type:    "otp",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if otpLoginResp.Login != nil && otpLoginResp.Login.IsOn {
		var loginParams = maps.Map{}
		err = json.Unmarshal(otpLoginResp.Login.ParamsJSON, &loginParams)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var secret = loginParams.GetString("secret")
		if gotp.NewDefaultTOTP(secret).Now() != params.OtpCode {
			this.FailField("otpCode", "请输入正确的OTP动态密码")
			return
		}
	}

	// 写入SESSION
	var localSid = rands.HexString(32)
	this.Data["localSid"] = localSid
	this.Data["ip"] = loginutils.RemoteIP(&this.ActionObject)
	params.Auth.StoreAdmin(adminId, params.Remember, localSid)

	// 删除OTP SESSION
	_, err = this.RPC().LoginSessionRPC().DeleteLoginSession(this.AdminContext(), &pb.DeleteLoginSessionRequest{Sid: sid})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 记录日志
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	err = dao.SharedLogDAO.CreateAdminLog(rpcClient.Context(adminId), oplogs.LevelInfo, this.Request.URL.Path, this.Lang(codes.AdminLogin_LogOtpVerifiedSuccess), loginutils.RemoteIP(&this.ActionObject), codes.AdminLogin_LogOtpVerifiedSuccess, nil)
	if err != nil {
		utils.PrintError(err)
	}

	this.Success()
}
