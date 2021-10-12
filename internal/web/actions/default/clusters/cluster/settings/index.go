package settings

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("basic")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	// 基本信息
	clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cluster := clusterResp.NodeCluster
	if cluster == nil {
		this.WriteString("not found cluster")
		return
	}

	// 认证
	var grantMap interface{} = nil

	if cluster.NodeGrantId > 0 {
		grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: cluster.NodeGrantId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		grant := grantResp.NodeGrant
		if grant != nil {
			grantMap = maps.Map{
				"id":         grant.Id,
				"name":       grant.Name,
				"method":     grant.Method,
				"methodName": grantutils.FindGrantMethodName(grant.Method),
			}
		}
	}
	this.Data["grant"] = grantMap

	// 时区
	this.Data["timeZoneGroups"] = nodeconfigs.FindAllTimeZoneGroups()
	this.Data["timeZoneLocations"] = nodeconfigs.FindAllTimeZoneLocations()

	if len(cluster.TimeZone) == 0 {
		cluster.TimeZone = nodeconfigs.DefaultTimeZoneLocation
	}
	this.Data["timeZoneLocation"] = nodeconfigs.FindTimeZoneLocation(cluster.TimeZone)

	this.Data["cluster"] = maps.Map{
		"id":         cluster.Id,
		"name":       cluster.Name,
		"installDir": cluster.InstallDir,
		"timeZone":   cluster.TimeZone,
	}

	this.Show()
}

// RunPost 保存设置
func (this *IndexAction) RunPost(params struct {
	ClusterId  int64
	Name       string
	GrantId    int64
	InstallDir string
	TimeZone   string

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改集群基础设置 %d", params.ClusterId)

	params.Must.
		Field("name", params.Name).
		Require("请输入集群名称")

	_, err := this.RPC().NodeClusterRPC().UpdateNodeCluster(this.AdminContext(), &pb.UpdateNodeClusterRequest{
		NodeClusterId: params.ClusterId,
		Name:          params.Name,
		NodeGrantId:   params.GrantId,
		InstallDir:    params.InstallDir,
		TimeZone:      params.TimeZone,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
