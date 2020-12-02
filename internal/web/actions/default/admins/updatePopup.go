package admins

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
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
	AdminId int64
}) {

	adminResp, err := this.RPC().AdminRPC().FindEnabledAdmin(this.AdminContext(), &pb.FindEnabledAdminRequest{AdminId: params.AdminId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	admin := adminResp.Admin

	this.Data["admin"] = maps.Map{
		"id":       admin.Id,
		"fullname": admin.Fullname,
		"username": admin.Username,
	}

	moduleMaps := configloaders.AllModuleMaps()
	for _, m := range moduleMaps {
		code := m.GetString("code")
		isChecked := false
		for _, module := range admin.Modules {
			if module.Code == code {
				isChecked = true
				break
			}
		}
		m["isChecked"] = isChecked
	}
	this.Data["modules"] = moduleMaps

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	AdminId int64

	Fullname    string
	Username    string
	Pass1       string
	Pass2       string
	ModuleCodes []string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改系统用户 %d", params.AdminId)

	params.Must.
		Field("fullname", params.Fullname).
		Require("请输入系统用户全名")

	params.Must.
		Field("username", params.Username).
		Require("请输入登录用户名").
		Match(`^[a-zA-Z0-9_]+$`, "用户名中只能包含英文、数字或下划线")

	existsResp, err := this.RPC().AdminRPC().CheckAdminUsername(this.AdminContext(), &pb.CheckAdminUsernameRequest{
		AdminId:  params.AdminId,
		Username: params.Username,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if existsResp.Exists {
		this.FailField("username", "此用户名已经被别的系统用户使用，请换一个")
	}

	if len(params.Pass1) > 0 {
		params.Must.
			Field("pass1", params.Pass1).
			Require("请输入登录密码").
			Field("pass2", params.Pass2).
			Require("请输入确认登录密码")
		if params.Pass1 != params.Pass2 {
			this.FailField("pass2", "两次输入的密码不一致")
		}
	}

	modules := []*systemconfigs.AdminModule{}
	for _, code := range params.ModuleCodes {
		modules = append(modules, &systemconfigs.AdminModule{
			Code:     code,
			AllowAll: true,
			Actions:  nil, // TODO 后期再开放细粒度控制
		})
	}
	modulesJSON, err := json.Marshal(modules)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().AdminRPC().UpdateAdmin(this.AdminContext(), &pb.UpdateAdminRequest{
		AdminId:     params.AdminId,
		Username:    params.Username,
		Password:    params.Pass1,
		Fullname:    params.Fullname,
		ModulesJSON: modulesJSON,
		IsSuper:     false, // TODO 后期再开放创建超级用户
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
