package login

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
		"username": admin.Username,
		"fullname": admin.Fullname,
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	Username  string
	Password  string
	Password2 string

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.AdminLogin_LogUpdateLogin)

	params.Must.
		Field("username", params.Username).
		Require("请输入登录用户名").
		Match(`^[a-zA-Z0-9_]+$`, "用户名中只能包含英文、数字或下划线")

	existsResp, err := this.RPC().AdminRPC().CheckAdminUsername(this.AdminContext(), &pb.CheckAdminUsernameRequest{
		AdminId:  this.AdminId(),
		Username: params.Username,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if existsResp.Exists {
		this.FailField("username", "此用户名已经被别的管理员使用，请换一个")
	}

	if len(params.Password) > 0 {
		if params.Password != params.Password2 {
			this.FailField("password2", "两次输入的密码不一致")
		}
	}

	_, err = this.RPC().AdminRPC().UpdateAdminLogin(this.AdminContext(), &pb.UpdateAdminLoginRequest{
		AdminId:  this.AdminId(),
		Username: params.Username,
		Password: params.Password,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 通知更新
	err = configloaders.NotifyAdminModuleMappingChange()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
