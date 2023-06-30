package accesskeys

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	AccessKeyId int64
}) {
	defer this.CreateLogInfo(codes.UserAccessKey_LogDeleteUserAccessKey, params.AccessKeyId)

	_, err := this.RPC().UserAccessKeyRPC().DeleteUserAccessKey(this.AdminContext(), &pb.DeleteUserAccessKeyRequest{UserAccessKeyId: params.AccessKeyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
