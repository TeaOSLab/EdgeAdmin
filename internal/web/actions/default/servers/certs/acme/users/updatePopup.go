package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	UserId int64
}) {
	userResp, err := this.RPC().ACMEUserRPC().FindEnabledACMEUser(this.AdminContext(), &pb.FindEnabledACMEUserRequest{AcmeUserId: params.UserId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	user := userResp.AcmeUser
	if user == nil {
		this.NotFound("acmeUser", params.UserId)
		return
	}

	// 服务商
	var providerMap maps.Map
	if user.AcmeProvider != nil {
		providerMap = maps.Map{
			"name": user.AcmeProvider.Name,
			"code": user.AcmeProvider.Code,
		}
	}

	// 账号
	var accountMap maps.Map
	if user.AcmeProviderAccount != nil {
		accountMap = maps.Map{
			"id":   user.AcmeProviderAccount.Id,
			"name": user.AcmeProviderAccount.Name,
		}
	}

	this.Data["user"] = maps.Map{
		"id":          user.Id,
		"email":       user.Email,
		"description": user.Description,
		"provider":    providerMap,
		"account":     accountMap,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	UserId      int64
	Description string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.ACMEUser_LogUpdateACMEUser, params.UserId)

	_, err := this.RPC().ACMEUserRPC().UpdateACMEUser(this.AdminContext(), &pb.UpdateACMEUserRequest{
		AcmeUserId:  params.UserId,
		Description: params.Description,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
