package nodes

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
	// 所有集群
	/**clusters, err := models.SharedNodeClusterDAO.FindAllEnableClusters()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range clusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps**/

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name      string
	ClusterId int

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称")

	// TODO 检查cluster
	if params.ClusterId <= 0 {
		this.Fail("请选择所在集群")
	}

	// TODO 检查SSH授权
}
