package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "node", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
	// 所有集群
	resp, err := this.RPC().NodeClusterRPC().FindAllEnabledClusters(this.AdminContext(), &pb.FindAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range resp.Clusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name      string
	ClusterId int64
	GrantId   int64
	SshHost   string
	SshPort   int

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称")

	// TODO 检查cluster
	if params.ClusterId <= 0 {
		this.Fail("请选择所在集群")
	}

	// TODO 检查登录授权
	loginInfo := &pb.NodeLogin{
		Id:   0,
		Name: "SSH",
		Type: "ssh",
		Params: maps.Map{
			"grantId": params.GrantId,
			"host":    params.SshHost,
			"port":    params.SshPort,
		}.AsJSON(),
	}

	// 保存
	_, err := this.RPC().NodeRPC().CreateNode(this.AdminContext(), &pb.CreateNodeRequest{
		Name:      params.Name,
		ClusterId: params.ClusterId,
		Login:     loginInfo,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
