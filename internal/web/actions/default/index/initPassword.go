// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package index

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type InitPasswordAction struct {
	actionutils.ParentAction
}

func (this *InitPasswordAction) Init() {
	this.Nav("", "", "")
}

func (this *InitPasswordAction) RunGet(params struct{}) {
	isNotInitialized, err := this.isNotInitialized()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if !isNotInitialized {
		this.RedirectURL("/")
		return
	}

	this.Data["username"] = "admin"
	this.Data["password"] = ""

	this.Show()
}

func (this *InitPasswordAction) RunPost(params struct {
	Username string
	Password string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	isNotInitialized, err := this.isNotInitialized()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if !isNotInitialized {
		this.ResponseWriter.WriteHeader(http.StatusForbidden)
		return
	}

	params.Must.
		Field("username", params.Username).
		Require("请输入登录用户名").
		Match(`^[a-zA-Z0-9_]+$`, "用户名中只能包含英文、数字或下划线").
		Field("password", params.Password).
		Require("请输入密码")

	// 查找ID
	adminResp, err := this.RPC().AdminRPC().FindAdminWithUsername(this.AdminContext(), &pb.FindAdminWithUsernameRequest{Username: "admin" /** 固定的 **/})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if adminResp.Admin == nil {
		this.Fail("数据错误，请将数据库中的edgeAdmins表中的用户名修改为admin后再试")
		return
	}
	var adminId = adminResp.Admin.Id

	// 修改密码
	_, err = this.RPC().AdminRPC().UpdateAdminLogin(this.AdminContext(), &pb.UpdateAdminLoginRequest{
		AdminId:  adminId,
		Username: params.Username,
		Password: params.Password, // raw
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 修改为初始化完成
	_, err = this.RPC().SysSettingRPC().UpdateSysSetting(this.AdminContext(), &pb.UpdateSysSettingRequest{
		Code:      systemconfigs.SettingCodeStandaloneInstanceInitialized,
		ValueJSON: []byte("1"),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}

func (this *InitPasswordAction) isNotInitialized() (bool, error) {
	settingResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeStandaloneInstanceInitialized})
	if err != nil {
		return false, err
	}
	return string(settingResp.ValueJSON) == "0", nil
}
