package db

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name     string
	Host     string
	Port     int32
	Database string
	Username string
	Password string

	Description string
	IsOn        bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称").
		Field("host", params.Host).
		Require("请输入主机地址").
		Field("port", params.Port).
		Gt(0, "请输入正确的数据库端口").
		Lt(65535, "请输入正确的数据库端口").
		Field("database", params.Database).
		Require("请输入数据库名称").
		Field("username", params.Username).
		Require("请输入连接数据库的用户名")

	createResp, err := this.RPC().DBNodeRPC().CreateDBNode(this.AdminContext(), &pb.CreateDBNodeRequest{
		IsOn:        params.IsOn,
		Name:        params.Name,
		Description: params.Description,
		Host:        params.Host,
		Port:        params.Port,
		Database:    params.Database,
		Username:    params.Username,
		Password:    params.Password,
		Charset:     "", // 暂时不能修改
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLogInfo(codes.DBNode_LogCreateDBNode, createResp.DbNodeId)

	this.Success()
}
