package settings

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

// IndexAction 服务基本信息设置
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("basic")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 所有集群
	resp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClusters(this.AdminContext(), &pb.FindAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var clusterMaps = []maps.Map{}
	for _, cluster := range resp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	// 当前服务信息
	serverResp, err := this.RPC().ServerRPC().FindEnabledServer(this.AdminContext(), &pb.FindEnabledServerRequest{
		ServerId:       params.ServerId,
		IgnoreSSLCerts: true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var server = serverResp.Server
	if server == nil {
		this.NotFound("server", params.ServerId)
		return
	}

	// 用户
	if server.User != nil {
		this.Data["user"] = maps.Map{
			"id":       server.User.Id,
			"fullname": server.User.Fullname,
			"username": server.User.Username,
		}
	} else {
		this.Data["user"] = nil
	}

	// 套餐
	this.initUserPlan(server)

	// 集群
	var clusterId = int64(0)
	this.Data["clusterName"] = ""
	if server.NodeCluster != nil {
		clusterId = server.NodeCluster.Id
		this.Data["clusterName"] = server.NodeCluster.Name
	}

	// 分组
	var groupMaps = []maps.Map{}
	if len(server.ServerGroups) > 0 {
		for _, group := range server.ServerGroups {
			groupMaps = append(groupMaps, maps.Map{
				"id":   group.Id,
				"name": group.Name,
			})
		}
	}

	// 域名和限流状态
	var trafficLimitStatus *serverconfigs.TrafficLimitStatus
	if len(server.Config) > 0 {
		var serverConfig = &serverconfigs.ServerConfig{}
		err = json.Unmarshal(server.Config, serverConfig)
		if err == nil {
			if serverConfig.TrafficLimitStatus != nil && serverConfig.TrafficLimitStatus.IsValid() {
				trafficLimitStatus = serverConfig.TrafficLimitStatus
			}
		}
	}

	this.Data["server"] = maps.Map{
		"id":                 server.Id,
		"clusterId":          clusterId,
		"type":               server.Type,
		"name":               server.Name,
		"description":        server.Description,
		"isOn":               server.IsOn,
		"groups":             groupMaps,
		"trafficLimitStatus": trafficLimitStatus,
	}

	var serverType = serverconfigs.FindServerType(server.Type)
	if serverType == nil {
		this.ErrorPage(errors.New("invalid server type '" + server.Type + "'"))
		return
	}

	var typeName = serverType.GetString("name")
	this.Data["typeName"] = typeName

	// 记录最近使用
	_, err = this.RPC().LatestItemRPC().IncreaseLatestItem(this.AdminContext(), &pb.IncreaseLatestItemRequest{
		ItemType: "server",
		ItemId:   params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}

// RunPost 保存
func (this *IndexAction) RunPost(params struct {
	ServerId       int64
	UserId         int64
	Name           string
	Description    string
	ClusterId      int64
	KeepOldConfigs bool
	GroupIds       []int64
	IsOn           bool
	UserPlanId     int64

	Must *actions.Must
}) {
	// 记录日志
	defer this.CreateLogInfo(codes.Server_LogUpdateServerBasic, params.ServerId)

	params.Must.
		Field("name", params.Name).
		Require("请输入服务名称")

	if params.ClusterId <= 0 {
		this.Fail("请选择部署的集群")
	}

	// 修改基本信息
	_, err := this.RPC().ServerRPC().UpdateServerBasic(this.AdminContext(), &pb.UpdateServerBasicRequest{
		ServerId:       params.ServerId,
		Name:           params.Name,
		Description:    params.Description,
		NodeClusterId:  params.ClusterId,
		KeepOldConfigs: params.KeepOldConfigs,
		IsOn:           params.IsOn,
		ServerGroupIds: params.GroupIds,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 修改用户
	if params.UserId > 0 {
		_, err = this.RPC().ServerRPC().UpdateServerUser(this.AdminContext(), &pb.UpdateServerUserRequest{
			ServerId: params.ServerId,
			UserId:   params.UserId,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		// 修改套餐
		if params.UserPlanId > 0 {
			_, err = this.RPC().ServerRPC().UpdateServerUserPlan(this.AdminContext(), &pb.UpdateServerUserPlanRequest{
				ServerId:   params.ServerId,
				UserPlanId: params.UserPlanId,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	}

	this.Success()
}
