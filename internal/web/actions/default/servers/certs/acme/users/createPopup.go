package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	PlatformUserId int64
	ProviderCode   string
}) {
	this.Data["platformUserId"] = params.PlatformUserId
	this.Data["providerCode"] = params.ProviderCode

	// 平台用户信息
	this.Data["platformUser"] = nil
	if params.PlatformUserId > 0 {
		platformUserResp, err := this.RPC().UserRPC().FindEnabledUser(this.AdminContext(), &pb.FindEnabledUserRequest{UserId: params.PlatformUserId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var platformUser = platformUserResp.User
		if platformUser != nil {
			this.Data["platformUser"] = maps.Map{
				"id":       platformUser.Id,
				"username": platformUser.Username,
				"fullname": platformUser.Fullname,
			}
		}
	}

	// 服务商
	providersResp, err := this.RPC().ACMEProviderRPC().FindAllACMEProviders(this.AdminContext(), &pb.FindAllACMEProvidersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var providerMaps = []maps.Map{}
	for _, provider := range providersResp.AcmeProviders {
		providerMaps = append(providerMaps, maps.Map{
			"code":       provider.Code,
			"name":       provider.Name,
			"requireEAB": provider.RequireEAB,
		})
	}
	this.Data["providers"] = providerMaps

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	PlatformUserId int64
	Email          string
	ProviderCode   string
	AccountId      int64
	Description    string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("email", params.Email).
		Require("请输入邮箱").
		Email("请输入正确的邮箱格式").
		Field("providerCode", params.ProviderCode).
		Require("请选择所属服务商")

	providerResp, err := this.RPC().ACMEProviderRPC().FindACMEProviderWithCode(this.AdminContext(), &pb.FindACMEProviderWithCodeRequest{
		AcmeProviderCode: params.ProviderCode,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if providerResp.AcmeProvider == nil {
		this.Fail("找不到要选择的证书")
	}
	if providerResp.AcmeProvider.RequireEAB {
		if params.AccountId <= 0 {
			this.Fail("此服务商要求必须选择或创建服务商账号")
		}

		// 同一个账号只能有一个用户
		countResp, err := this.RPC().ACMEUserRPC().
			CountACMEUsers(this.AdminContext(), &pb.CountAcmeUsersRequest{
				AcmeProviderAccountId: params.AccountId,
			})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if countResp.Count > 0 {
			this.Fail("此服务商账号已被别的用户使用，请换成别的账号")
		}
	}

	createResp, err := this.RPC().ACMEUserRPC().CreateACMEUser(this.AdminContext(), &pb.CreateACMEUserRequest{
		UserId:                params.PlatformUserId,
		Email:                 params.Email,
		Description:           params.Description,
		AcmeProviderCode:      params.ProviderCode,
		AcmeProviderAccountId: params.AccountId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 返回数据
	this.Data["acmeUser"] = maps.Map{
		"id":           createResp.AcmeUserId,
		"description":  params.Description,
		"email":        params.Email,
		"providerCode": params.ProviderCode,
	}

	// 日志
	defer this.CreateLogInfo(codes.ACMEUser_LogCreateACMEUser, createResp.AcmeUserId)

	this.Success()
}
