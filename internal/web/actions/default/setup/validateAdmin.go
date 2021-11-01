package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type ValidateAdminAction struct {
	actionutils.ParentAction
}

func (this *ValidateAdminAction) RunPost(params struct {
	AdminUsername  string
	AdminPassword  string
	AdminPassword2 string
	Must           *actions.Must
}) {
	params.Must.
		Field("adminUsername", params.AdminUsername).
		Require("请输入管理员登录用户名").
		Match(`^[a-zA-Z0-9_]+$`, "用户名中只能包含英文、数字或下划线").
		Field("adminPassword", params.AdminPassword).
		Require("请输入管理员登录密码").
		Match(`^[a-zA-Z0-9_]+$`, "密码中只能包含英文、数字或下划线").
		Field("adminPassword2", params.AdminPassword2).
		Require("请输入确认密码").
		Equal(params.AdminPassword, "两次输入的密码不一致")

	this.Data["admin"] = maps.Map{
		"username":     params.AdminUsername,
		"password":     params.AdminPassword,
		"passwordMask": strings.Repeat("*", len(params.AdminPassword)),
	}

	this.Success()
}
