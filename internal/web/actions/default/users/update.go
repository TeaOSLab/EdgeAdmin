package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/users/userutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
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
	UserId int64
}) {
	err := userutils.InitUser(this.Parent(), params.UserId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

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

	// AccessKey数量
	countAccessKeyResp, err := this.RPC().UserAccessKeyRPC().CountAllEnabledUserAccessKeys(this.AdminContext(), &pb.CountAllEnabledUserAccessKeysRequest{UserId: params.UserId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var countAccessKeys = countAccessKeyResp.Count

	// 是否有实名认证
	hasNewIndividualIdentity, hasNewEnterpriseIdentity, identityTag, err := userutils.CheckUserIdentity(this.RPC(), this.AdminContext(), params.UserId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// OTP认证
	var otpLoginIsOn = false
	if user.OtpLogin != nil {
		otpLoginIsOn = user.OtpLogin.IsOn
	}

	this.Data["user"] = maps.Map{
		"id":              user.Id,
		"username":        user.Username,
		"fullname":        user.Fullname,
		"email":           user.Email,
		"tel":             user.Tel,
		"remark":          user.Remark,
		"mobile":          user.Mobile,
		"isOn":            user.IsOn,
		"countAccessKeys": countAccessKeys,
		"bandwidthAlgo":   user.BandwidthAlgo,

		// 实名认证
		"hasNewIndividualIdentity": hasNewIndividualIdentity,
		"hasNewEnterpriseIdentity": hasNewEnterpriseIdentity,
		"identityTag":              identityTag,

		// otp
		"otpLoginIsOn": otpLoginIsOn,
	}

	this.Data["clusterId"] = 0
	if user.NodeCluster != nil {
		this.Data["clusterId"] = user.NodeCluster.Id
	}

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	UserId        int64
	Username      string
	Pass1         string
	Pass2         string
	Fullname      string
	Mobile        string
	Tel           string
	Email         string
	Remark        string
	IsOn          bool
	ClusterId     int64
	BandwidthAlgo string

	// OTP
	OtpOn bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.User_LogUpdateUser, params.UserId)

	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Match(`^[a-zA-Z0-9_]+$`, "用户名中只能含有英文、数字和下划线")

	checkUsernameResp, err := this.RPC().UserRPC().CheckUserUsername(this.AdminContext(), &pb.CheckUserUsernameRequest{
		UserId:   params.UserId,
		Username: params.Username,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if checkUsernameResp.Exists {
		this.FailField("username", "此用户名已经被占用，请换一个")
	}

	if len(params.Pass1) > 0 {
		params.Must.
			Field("pass1", params.Pass1).
			Require("请输入密码").
			Field("pass2", params.Pass2).
			Require("请再次输入确认密码").
			Equal(params.Pass1, "两次输入的密码不一致")
	}

	params.Must.
		Field("fullname", params.Fullname).
		Require("请输入全名")

	if len(params.Mobile) > 0 {
		params.Must.
			Field("mobile", params.Mobile).
			Mobile("请输入正确的手机号")
	}
	if len(params.Email) > 0 {
		params.Must.
			Field("email", params.Email).
			Email("请输入正确的电子邮箱")
	}

	_, err = this.RPC().UserRPC().UpdateUser(this.AdminContext(), &pb.UpdateUserRequest{
		UserId:        params.UserId,
		Username:      params.Username,
		Password:      params.Pass1,
		Fullname:      params.Fullname,
		Mobile:        params.Mobile,
		Tel:           params.Tel,
		Email:         params.Email,
		Remark:        params.Remark,
		IsOn:          params.IsOn,
		NodeClusterId: params.ClusterId,
		BandwidthAlgo: params.BandwidthAlgo,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 修改OTP
	otpLoginResp, err := this.RPC().LoginRPC().FindEnabledLogin(this.AdminContext(), &pb.FindEnabledLoginRequest{
		UserId: params.UserId,
		Type:   "otp",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	{
		var otpLogin = otpLoginResp.Login
		if params.OtpOn {
			if otpLogin == nil {
				otpLogin = &pb.Login{
					Id:   0,
					Type: "otp",
					ParamsJSON: maps.Map{
						"secret": gotp.RandomSecret(16), // TODO 改成可以设置secret长度
					}.AsJSON(),
					IsOn:   true,
					UserId: params.UserId,
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
				Type:   "otp",
				IsOn:   false,
				UserId: params.UserId,
			}})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	}

	this.Success()
}
