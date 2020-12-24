package admins

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
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
	AdminId int64
}) {
	loginResp, err := this.RPC().LoginRPC().FindEnabledLogin(this.AdminContext(), &pb.FindEnabledLoginRequest{
		AdminId: params.AdminId,
		Type:    "otp",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	login := loginResp.Login
	if login == nil || !login.IsOn {
		this.NotFound("adminLogin", params.AdminId)
		return
	}

	loginParams := maps.Map{}
	err = json.Unmarshal(login.ParamsJSON, &loginParams)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	secret := loginParams.GetString("secret")

	// 当前用户信息
	adminResp, err := this.RPC().AdminRPC().FindEnabledAdmin(this.AdminContext(), &pb.FindEnabledAdminRequest{AdminId: params.AdminId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	admin := adminResp.Admin
	if admin == nil {
		this.NotFound("admin", params.AdminId)
		return
	}

	uiConfig, err := configloaders.LoadAdminUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	url := gotp.NewDefaultTOTP(secret).ProvisioningUri(admin.Username, uiConfig.AdminSystemName)
	data, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.AddHeader("Content-Type", "image/png")
	this.Write(data)
}
