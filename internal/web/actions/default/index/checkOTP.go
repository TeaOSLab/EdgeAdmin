package index

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

// 检查是否需要OTP
type CheckOTPAction struct {
	actionutils.ParentAction
}

func (this *CheckOTPAction) Init() {
	this.Nav("", "", "")
}

func (this *CheckOTPAction) RunPost(params struct {
	Username string

	Must *actions.Must
}) {
	checkResp, err := this.RPC().AdminRPC().CheckAdminOTPWithUsername(this.AdminContext(), &pb.CheckAdminOTPWithUsernameRequest{Username: params.Username})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["requireOTP"] = checkResp.RequireOTP
	this.Success()
}
