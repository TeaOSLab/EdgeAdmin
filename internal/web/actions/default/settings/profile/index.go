package profile

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	adminResp, err := this.RPC().AdminRPC().FindEnabledAdmin(this.AdminContext(), &pb.FindEnabledAdminRequest{AdminId: this.AdminId()})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	admin := adminResp.Admin
	if admin == nil {
		this.NotFound("admin", this.AdminId())
		return
	}

	this.Data["admin"] = maps.Map{
		"fullname": admin.Fullname,
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	Fullname string

	Must *actions.Must
}) {
	defer this.CreateLogInfo("修改个人资料")

	params.Must.
		Field("fullname", params.Fullname).
		Require("请输入你的姓名")

	_, err := this.RPC().AdminRPC().UpdateAdmin(this.AdminContext(), &pb.UpdateAdminRequest{
		AdminId:  this.AdminId(),
		Fullname: params.Fullname,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
