package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs/nodes"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "", "create")
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

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入服务名称")

	if params.ClusterId <= 0 {
		this.Fail("请选择部署的集群")
	}

	// TODO 验证集群ID

	// 配置
	serverConfig := &nodes.ServerConfig{}
	serverConfig.IsOn = true
	serverConfig.Name = params.Name
	serverConfigJSON, err := serverConfig.AsJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 包含条件
	includeNodes := []maps.Map{}
	includeNodesJSON, err := json.Marshal(includeNodes)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 排除条件
	excludeNodes := []maps.Map{}
	excludeNodesJSON, err := json.Marshal(excludeNodes)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ServerRPC().CreateServer(this.AdminContext(), &pb.CreateServerRequest{
		UserId:           0,
		AdminId:          this.AdminId(),
		ClusterId:        params.ClusterId,
		Config:           serverConfigJSON,
		IncludeNodesJSON: includeNodesJSON,
		ExcludeNodesJSON: excludeNodesJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
