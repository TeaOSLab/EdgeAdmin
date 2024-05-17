// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreateNodeInstallAction struct {
	actionutils.ParentAction
}

func (this *CreateNodeInstallAction) RunPost(params struct {
	NodeId  int64
	SshHost string
	SshPort int
	GrantId int64

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.NodeSSH_LogUpdateNodeSSH, params.NodeId)

	params.Must.
		Field("sshHost2", params.SshHost).
		Require("请填写SSH主机地址").
		Field("sshPort2", params.SshPort).
		Gt(0, "请填写SSH主机端口").
		Lt(65535, "SSH主机端口需要小于65535").
		Field("grantId", params.GrantId).
		Gt(0, "请选择SSH登录认证")

	// 查询login
	nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var node = nodeResp.Node
	if node == nil {
		this.Fail("找不到要修改的节点")
	}
	var loginId int64
	if node.NodeLogin != nil {
		loginId = node.NodeLogin.Id
	}

	// 修改节点信息
	_, err = this.RPC().NodeRPC().UpdateNodeLogin(this.AdminContext(), &pb.UpdateNodeLoginRequest{
		NodeId: params.NodeId,
		NodeLogin: &pb.NodeLogin{
			Id:   loginId,
			Name: "SSH",
			Type: "ssh",
			Params: maps.Map{
				"grantId": params.GrantId,
				"host":    params.SshHost,
				"port":    params.SshPort,
			}.AsJSON(),
		},
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 开始安装
	_, err = this.RPC().NodeRPC().InstallNode(this.AdminContext(), &pb.InstallNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.Fail("安装失败：" + err.Error())
	}

	this.Success()
}
