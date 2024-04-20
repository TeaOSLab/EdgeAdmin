package admins

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/otputils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

type OtpQrcodeAction struct {
	actionutils.ParentAction
}

func (this *OtpQrcodeAction) Init() {
	this.Nav("", "", "")
}

func (this *OtpQrcodeAction) RunGet(params struct {
	AdminId  int64
	Download bool
}) {
	loginResp, err := this.RPC().LoginRPC().FindEnabledLogin(this.AdminContext(), &pb.FindEnabledLoginRequest{
		AdminId: params.AdminId,
		Type:    "otp",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var login = loginResp.Login
	if login == nil || !login.IsOn {
		this.NotFound("adminLogin", params.AdminId)
		return
	}

	var loginParams = maps.Map{}
	err = json.Unmarshal(login.ParamsJSON, &loginParams)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var secret = loginParams.GetString("secret")

	// 当前用户信息
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

	uiConfig, err := configloaders.LoadAdminUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var url = gotp.NewDefaultTOTP(secret).ProvisioningUri(admin.Username, uiConfig.AdminSystemName)

	data, err := qrcode.Encode(otputils.FixIssuer(url), qrcode.Medium, 256)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if params.Download {
		var filename = "OTP-ADMIN-" + admin.Username + ".png"
		this.AddHeader("Content-Disposition", "attachment; filename=\""+filename+"\";")
	}
	this.AddHeader("Content-Type", "image/png")
	this.AddHeader("Content-Length", types.String(len(data)))
	_, _ = this.Write(data)
}
