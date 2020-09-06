package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "cluster", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name       string
	GrantId    int64
	InstallDir string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入集群名称")

	_, err := this.RPC().NodeClusterRPC().CreateNodeCluster(this.AdminContext(), &pb.CreateNodeClusterRequest{
		Name:       params.Name,
		GrantId:    params.GrantId,
		InstallDir: params.InstallDir,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
