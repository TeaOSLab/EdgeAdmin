package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	NodeId int64
}) {
	nodeResp, err := this.RPC().AuthorityNodeRPC().FindEnabledAuthorityNode(this.AdminContext(), &pb.FindEnabledAuthorityNodeRequest{
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
		"isOn":        node.IsOn,
	}

	this.Show()
}

// 保存基础设置
func (this *UpdateAction) RunPost(params struct {
	NodeId      int64
	Name        string
	Description string
	IsOn        bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入认证节点名称")

	_, err := this.RPC().AuthorityNodeRPC().UpdateAuthorityNode(this.AdminContext(), &pb.UpdateAuthorityNodeRequest{
		NodeId:      params.NodeId,
		Name:        params.Name,
		Description: params.Description,
		IsOn:        params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改认证节点 %d", params.NodeId)

	this.Success()
}
