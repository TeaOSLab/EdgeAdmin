package node

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/ipAddresses/ipaddressutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "node", "update")
	this.SecondMenu("basic")
}

func (this *UpdateAction) RunGet(params struct {
	NodeId int64
}) {
	_, err := nodeutils.InitNodeInfo(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["nodeId"] = params.NodeId

	nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	node := nodeResp.Node
	if node == nil {
		this.WriteString("找不到要操作的节点")
		return
	}

	var clusterMap maps.Map = nil
	if node.NodeCluster != nil {
		clusterMap = maps.Map{
			"id":   node.NodeCluster.Id,
			"name": node.NodeCluster.Name,
		}
	}

	// IP地址
	ipAddressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledNodeIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledNodeIPAddressesWithNodeIdRequest{
		NodeId: params.NodeId,
		Role:   nodeconfigs.NodeRoleNode,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	ipAddressMaps := []maps.Map{}
	for _, addr := range ipAddressesResp.NodeIPAddresses {
		thresholds, err := ipaddressutils.InitNodeIPAddressThresholds(this.Parent(), addr.Id)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		ipAddressMaps = append(ipAddressMaps, maps.Map{
			"id":         addr.Id,
			"name":       addr.Name,
			"ip":         addr.Ip,
			"canAccess":  addr.CanAccess,
			"isOn":       addr.IsOn,
			"isUp":       addr.IsUp,
			"thresholds": thresholds,
		})
	}

	// 分组
	var groupMap maps.Map = nil
	if node.NodeGroup != nil {
		groupMap = maps.Map{
			"id":   node.NodeGroup.Id,
			"name": node.NodeGroup.Name,
		}
	}

	// 区域
	var regionMap maps.Map = nil
	if node.NodeRegion != nil {
		regionMap = maps.Map{
			"id":   node.NodeRegion.Id,
			"name": node.NodeRegion.Name,
		}
	}

	var m = maps.Map{
		"id":          node.Id,
		"name":        node.Name,
		"ipAddresses": ipAddressMaps,
		"cluster":     clusterMap,
		"isOn":        node.IsOn,
		"group":       groupMap,
		"region":      regionMap,
	}

	if node.NodeCluster != nil {
		m["primaryCluster"] = maps.Map{
			"id":   node.NodeCluster.Id,
			"name": node.NodeCluster.Name,
		}
	} else {
		m["primaryCluster"] = nil
	}

	if len(node.SecondaryNodeClusters) > 0 {
		var secondaryClusterMaps = []maps.Map{}
		for _, cluster := range node.SecondaryNodeClusters {
			secondaryClusterMaps = append(secondaryClusterMaps, maps.Map{
				"id":   cluster.Id,
				"name": cluster.Name,
			})
		}
		m["secondaryClusters"] = secondaryClusterMaps
	} else {
		m["secondaryClusters"] = []interface{}{}
	}

	this.Data["node"] = m

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	LoginId             int64
	NodeId              int64
	GroupId             int64
	RegionId            int64
	Name                string
	IPAddressesJSON     []byte `alias:"ipAddressesJSON"`
	PrimaryClusterId    int64
	SecondaryClusterIds []byte
	IsOn                bool

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改节点 %d", params.NodeId)

	if params.NodeId <= 0 {
		this.Fail("要操作的节点不存在")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称")

	// TODO 检查cluster
	if params.PrimaryClusterId <= 0 {
		this.Fail("请选择节点所在主集群")
	}

	var secondaryClusterIds = []int64{}
	if len(params.SecondaryClusterIds) > 0 {
		err := json.Unmarshal(params.SecondaryClusterIds, &secondaryClusterIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// IP地址
	ipAddresses := []maps.Map{}
	if len(params.IPAddressesJSON) > 0 {
		err := json.Unmarshal(params.IPAddressesJSON, &ipAddresses)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	if len(ipAddresses) == 0 {
		this.Fail("请至少输入一个IP地址")
	}

	// 保存
	_, err := this.RPC().NodeRPC().UpdateNode(this.AdminContext(), &pb.UpdateNodeRequest{
		NodeId:                  params.NodeId,
		NodeGroupId:             params.GroupId,
		NodeRegionId:            params.RegionId,
		Name:                    params.Name,
		NodeClusterId:           params.PrimaryClusterId,
		SecondaryNodeClusterIds: secondaryClusterIds,
		IsOn:                    params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 禁用老的IP地址
	_, err = this.RPC().NodeIPAddressRPC().DisableAllNodeIPAddressesWithNodeId(this.AdminContext(), &pb.DisableAllNodeIPAddressesWithNodeIdRequest{
		NodeId: params.NodeId,
		Role:   nodeconfigs.NodeRoleNode,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 添加新的IP地址
	err = ipaddressutils.UpdateNodeIPAddresses(this.Parent(), params.NodeId, nodeconfigs.NodeRoleNode, params.IPAddressesJSON)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
