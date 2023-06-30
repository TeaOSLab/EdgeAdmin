package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	UserId int64
}) {
	defer this.CreateLogInfo(codes.User_LogDeleteUser, params.UserId)

	// TODO 检查用户是否有未完成的业务

	_, err := this.RPC().UserRPC().DeleteUser(this.AdminContext(), &pb.DeleteUserRequest{UserId: params.UserId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
