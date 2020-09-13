package node

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "node", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreateAction) RunPost(params struct {
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

	_, err := this.RPC().APINodeRPC().CreateAPINode(this.AdminContext(), &pb.CreateAPINodeRequest{
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
