package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "user")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().ACMEUserRPC().CountACMEUsers(this.AdminContext(), &pb.CountAcmeUsersRequest{
		AdminId: this.AdminId(),
		UserId:  0,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	usersResp, err := this.RPC().ACMEUserRPC().ListACMEUsers(this.AdminContext(), &pb.ListACMEUsersRequest{
		AdminId: this.AdminId(),
		UserId:  0,
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	userMaps := []maps.Map{}
	for _, user := range usersResp.AcmeUsers {
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

		userMaps = append(userMaps, maps.Map{
			"id":          user.Id,
			"email":       user.Email,
			"description": user.Description,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", user.CreatedAt),
			"provider":    providerMap,
			"account":     accountMap,
		})
	}
	this.Data["users"] = userMaps

	this.Show()
}
