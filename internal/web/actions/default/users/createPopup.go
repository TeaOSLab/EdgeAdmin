package users

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/userconfigs"
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
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Username     string
	Pass1        string
	Pass2        string
	Fullname     string
	Mobile       string
	Tel          string
	Email        string
	Remark       string
	ClusterId    int64
	FeaturesType string

	// OTP
	OtpOn bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var userId int64

	defer func() {
		this.CreateLogInfo(codes.User_LogCreateUser, userId)
	}()

	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Match(`^[a-zA-Z0-9_]+$`, "用户名中只能含有英文、数字和下划线")

	checkUsernameResp, err := this.RPC().UserRPC().CheckUserUsername(this.AdminContext(), &pb.CheckUserUsernameRequest{
		UserId:   0,
		Username: params.Username,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if checkUsernameResp.Exists {
		this.FailField("username", "此用户名已经被占用，请换一个")
	}

	params.Must.
		Field("pass1", params.Pass1).
		Require("请输入密码").
		Field("pass2", params.Pass2).
		Require("请再次输入确认密码").
		Equal(params.Pass1, "两次输入的密码不一致")

	params.Must.
		Field("fullname", params.Fullname).
		Require("请输入全名")

	if params.ClusterId <= 0 {
		this.Fail("请选择关联集群")
	}

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

	createResp, err := this.RPC().UserRPC().CreateUser(this.AdminContext(), &pb.CreateUserRequest{
		Username:      params.Username,
		Password:      params.Pass1,
		Fullname:      params.Fullname,
		Mobile:        params.Mobile,
		Tel:           params.Tel,
		Email:         params.Email,
		Remark:        params.Remark,
		Source:        "admin:" + numberutils.FormatInt64(this.AdminId()),
		NodeClusterId: params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	userId = createResp.UserId

	// 功能
	if teaconst.IsPlus {
		if params.FeaturesType == "default" {
			resp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeUserRegisterConfig})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			var config = userconfigs.DefaultUserRegisterConfig()
			if len(resp.ValueJSON) > 0 {
				err = json.Unmarshal(resp.ValueJSON, config)
				if err != nil {
					this.ErrorPage(err)
					return
				}
				_, err = this.RPC().UserRPC().UpdateUserFeatures(this.AdminContext(), &pb.UpdateUserFeaturesRequest{
					UserId:       userId,
					FeatureCodes: config.Features,
				})
				if err != nil {
					this.ErrorPage(err)
					return
				}
			}
		} else if params.FeaturesType == "all" {
			featuresResp, err := this.RPC().UserRPC().FindAllUserFeatureDefinitions(this.AdminContext(), &pb.FindAllUserFeatureDefinitionsRequest{})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			var featureCodes = []string{}
			for _, def := range featuresResp.Features {
				featureCodes = append(featureCodes, def.Code)
			}
			_, err = this.RPC().UserRPC().UpdateUserFeatures(this.AdminContext(), &pb.UpdateUserFeaturesRequest{
				UserId:       userId,
				FeatureCodes: featureCodes,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
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
			AdminId: 0,
			UserId:  userId,
		}})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Success()
}
