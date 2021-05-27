package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type OptionsAction struct {
	actionutils.ParentAction
}

func (this *OptionsAction) RunPost(params struct{}) {
	usersResp, err := this.RPC().UserRPC().ListEnabledUsers(this.AdminContext(), &pb.ListEnabledUsersRequest{
		Offset: 0,
		Size:   10000, // TODO 改进 <ns-user-selector> 组件
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	userMaps := []maps.Map{}
	for _, user := range usersResp.Users {
		userMaps = append(userMaps, maps.Map{
			"id":       user.Id,
			"fullname": user.Fullname,
			"username": user.Username,
		})
	}
	this.Data["users"] = userMaps

	this.Success()
}
