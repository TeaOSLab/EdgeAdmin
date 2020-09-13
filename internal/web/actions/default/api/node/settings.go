package node

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type SettingsAction struct {
	actionutils.ParentAction
}

func (this *SettingsAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("basic")
}

func (this *SettingsAction) RunGet(params struct {
	NodeId int64
}) {
	nodeResp, err := this.RPC().APINodeRPC().FindEnabledAPINode(this.AdminContext(), &pb.FindEnabledAPINodeRequest{
		NodeId: params.NodeId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	node := nodeResp.Node
	if node == nil {
		this.WriteString("要操作的节点不存在")
		return
	}

	this.Data["node"] = maps.Map{
		"id":          node.Id,
		"name":        node.Name,
		"description": node.Description,
		"host":        node.Host,
		"port":        node.Port,
	}

	this.Show()
}

// 保存基础设置
func (this *SettingsAction) RunPost(params struct {
	NodeId      int64
	Name        string
	Host        string
	Port        int
	Description string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入API节点").
		Field("host", params.Host).
		Require("请输入主机地址").
		Field("port", params.Port).
		Gt(0, "端口不能小于1").
		Lte(65535, "端口不能大于65535")

	_, err := this.RPC().APINodeRPC().UpdateAPINode(this.AdminContext(), &pb.UpdateAPINodeRequest{
		NodeId:      params.NodeId,
		Name:        params.Name,
		Description: params.Description,
		Host:        params.Host,
		Port:        int32(params.Port),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
