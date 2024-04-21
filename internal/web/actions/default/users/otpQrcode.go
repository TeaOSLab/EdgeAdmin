package users

import (
	"encoding/json"
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
	UserId   int64
	Download bool
}) {
	loginResp, err := this.RPC().LoginRPC().FindEnabledLogin(this.AdminContext(), &pb.FindEnabledLoginRequest{
		UserId: params.UserId,
		Type:   "otp",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var login = loginResp.Login
	if login == nil || !login.IsOn {
		this.NotFound("userLogin", params.UserId)
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
	userResp, err := this.RPC().UserRPC().FindEnabledUser(this.AdminContext(), &pb.FindEnabledUserRequest{UserId: params.UserId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var user = userResp.User
	if user == nil {
		this.NotFound("user", params.UserId)
		return
	}

	productName, err := this.findProductName()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var url = gotp.NewDefaultTOTP(secret).ProvisioningUri(user.Username, productName)
	data, err := qrcode.Encode(otputils.FixIssuer(url), qrcode.Medium, 256)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if params.Download {
		var filename = "OTP-USER-" + user.Username + ".png"
		this.AddHeader("Content-Disposition", "attachment; filename=\""+filename+"\";")
	}
	this.AddHeader("Content-Type", "image/png")
	this.AddHeader("Content-Length", types.String(len(data)))
	_, _ = this.Write(data)
}
