package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type OptionsAction struct {
	actionutils.ParentAction
}

func (this *OptionsAction) RunPost(params struct {
	Keyword string
}) {
	usersResp, err := this.RPC().UserRPC().ListEnabledUsers(this.AdminContext(), &pb.ListEnabledUsersRequest{
		Keyword: params.Keyword,
		Offset:  0,
		Size:    100,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var userMaps = []maps.Map{}
	for _, user := range usersResp.Users {
		userMaps = append(userMaps, maps.Map{
			"id":       user.Id,
			"fullname": user.Fullname,
			"username": user.Username,
			"name":     user.Fullname + "(" + user.Username + ")",
		})
	}
	this.Data["users"] = userMaps

	this.Success()
}
