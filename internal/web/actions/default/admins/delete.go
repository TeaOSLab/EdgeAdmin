package admins

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	AdminId int64
}) {
	defer this.CreateLogInfo(codes.Admin_LogDeleteAdmin, params.AdminId)

	_, err := this.RPC().AdminRPC().DeleteAdmin(this.AdminContext(), &pb.DeleteAdminRequest{AdminId: params.AdminId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 通知更改
	err = configloaders.NotifyAdminModuleMappingChange()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
