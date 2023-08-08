package ipaddressutils

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

// UpdateNodeIPAddresses 保存一组IP地址
func UpdateNodeIPAddresses(parentAction *actionutils.ParentAction, nodeId int64, role nodeconfigs.NodeRole, ipAddressesJSON []byte) error {
	var addresses = []maps.Map{}
	err := json.Unmarshal(ipAddressesJSON, &addresses)
	if err != nil {
		return err
	}
	for _, addr := range addresses {
		var resultAddrIds = []int64{}
		var addrId = addr.GetInt64("id")

		// 专属集群
		var addrClusterIds = []int64{}
		var addrClusters = addr.GetSlice("clusters")
		if len(addrClusters) > 0 {
			for _, addrCluster := range addrClusters {
				var m = maps.NewMap(addrCluster)
				var clusterId = m.GetInt64("id")
				if clusterId > 0 {
					addrClusterIds = append(addrClusterIds, clusterId)
				}
			}
		}

		if addrId > 0 {
			resultAddrIds = append(resultAddrIds, addrId)

			var isOn bool
			if !addr.Has("isOn") { // 兼容老版本
				isOn = true
			} else {
				isOn = addr.GetBool("isOn")
			}

			_, err = parentAction.RPC().NodeIPAddressRPC().UpdateNodeIPAddress(parentAction.AdminContext(), &pb.UpdateNodeIPAddressRequest{
				NodeIPAddressId: addrId,
				Ip:              addr.GetString("ip"),
				Name:            addr.GetString("name"),
				CanAccess:       addr.GetBool("canAccess"),
				IsOn:            isOn,
				IsUp:            addr.GetBool("isUp"),
				ClusterIds:      addrClusterIds,
			})
			if err != nil {
				return err
			}
		} else {
			var ipStrings = addr.GetString("ip")
			result, _ := utils.ExtractIP(ipStrings)

			if len(result) == 1 {
				// 单个创建
				createResp, err := parentAction.RPC().NodeIPAddressRPC().CreateNodeIPAddress(parentAction.AdminContext(), &pb.CreateNodeIPAddressRequest{
					NodeId:         nodeId,
					Role:           role,
					Name:           addr.GetString("name"),
					Ip:             result[0],
					CanAccess:      addr.GetBool("canAccess"),
					IsUp:           addr.GetBool("isUp"),
					NodeClusterIds: addrClusterIds,
				})
				if err != nil {
					return err
				}
				addrId = createResp.NodeIPAddressId
				resultAddrIds = append(resultAddrIds, addrId)
			} else if len(result) > 1 {
				// 批量创建
				createResp, err := parentAction.RPC().NodeIPAddressRPC().CreateNodeIPAddresses(parentAction.AdminContext(), &pb.CreateNodeIPAddressesRequest{
					NodeId:         nodeId,
					Role:           role,
					Name:           addr.GetString("name"),
					IpList:         result,
					CanAccess:      addr.GetBool("canAccess"),
					IsUp:           addr.GetBool("isUp"),
					GroupValue:     ipStrings,
					NodeClusterIds: addrClusterIds,
				})
				if err != nil {
					return err
				}
				resultAddrIds = append(resultAddrIds, createResp.NodeIPAddressIds...)
			}
		}

		// 保存阈值
		var thresholds = addr.GetSlice("thresholds")
		if len(thresholds) > 0 {
			thresholdsJSON, err := json.Marshal(thresholds)
			if err != nil {
				return err
			}

			for _, addrId := range resultAddrIds {
				_, err = parentAction.RPC().NodeIPAddressThresholdRPC().UpdateAllNodeIPAddressThresholds(parentAction.AdminContext(), &pb.UpdateAllNodeIPAddressThresholdsRequest{
					NodeIPAddressId:             addrId,
					NodeIPAddressThresholdsJSON: thresholdsJSON,
				})
				if err != nil {
					return err
				}
			}
		} else {
			for _, addrId := range resultAddrIds {
				_, err = parentAction.RPC().NodeIPAddressThresholdRPC().UpdateAllNodeIPAddressThresholds(parentAction.AdminContext(), &pb.UpdateAllNodeIPAddressThresholdsRequest{
					NodeIPAddressId:             addrId,
					NodeIPAddressThresholdsJSON: []byte("[]"),
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// InitNodeIPAddressThresholds 初始化IP阈值
func InitNodeIPAddressThresholds(parentAction *actionutils.ParentAction, addrId int64) ([]*nodeconfigs.IPAddressThresholdConfig, error) {
	thresholdsResp, err := parentAction.RPC().NodeIPAddressThresholdRPC().FindAllEnabledNodeIPAddressThresholds(parentAction.AdminContext(), &pb.FindAllEnabledNodeIPAddressThresholdsRequest{NodeIPAddressId: addrId})
	if err != nil {
		return nil, err
	}
	var thresholds = []*nodeconfigs.IPAddressThresholdConfig{}
	if len(thresholdsResp.NodeIPAddressThresholds) > 0 {
		for _, pbThreshold := range thresholdsResp.NodeIPAddressThresholds {
			var threshold = &nodeconfigs.IPAddressThresholdConfig{
				Id:      pbThreshold.Id,
				Items:   []*nodeconfigs.IPAddressThresholdItemConfig{},
				Actions: []*nodeconfigs.IPAddressThresholdActionConfig{},
			}
			if len(pbThreshold.ItemsJSON) > 0 {
				err = json.Unmarshal(pbThreshold.ItemsJSON, &threshold.Items)
				if err != nil {
					return nil, err
				}
			}
			if len(pbThreshold.ActionsJSON) > 0 {
				err = json.Unmarshal(pbThreshold.ActionsJSON, &threshold.Actions)
				if err != nil {
					return nil, err
				}
			}
			thresholds = append(thresholds, threshold)
		}
	}
	return thresholds, nil
}

// FindNodeClusterMapsWithNodeId 根据节点读取集群信息
func FindNodeClusterMapsWithNodeId(parentAction *actionutils.ParentAction, nodeId int64) ([]maps.Map, error) {
	var clusterMaps = []maps.Map{}
	if nodeId > 0 { // CDN边缘节点
		nodeResp, err := parentAction.RPC().NodeRPC().FindEnabledNode(parentAction.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: nodeId})
		if err != nil {
			return nil, err
		}
		var node = nodeResp.Node
		if node != nil {
			var clusters = []*pb.NodeCluster{}
			if node.NodeCluster != nil {
				clusters = append(clusters, nodeResp.Node.NodeCluster)
			}
			if len(node.SecondaryNodeClusters) > 0 {
				clusters = append(clusters, node.SecondaryNodeClusters...)
			}
			for _, cluster := range clusters {
				clusterMaps = append(clusterMaps, maps.Map{
					"id":        cluster.Id,
					"name":      cluster.Name,
					"isChecked": false,
				})
			}
		}
	}
	return clusterMaps, nil
}

// FindNodeClusterMaps 根据一组集群ID读取集群信息
func FindNodeClusterMaps(parentAction *actionutils.ParentAction, clusterIds []int64) ([]maps.Map, error) {
	var clusterMaps = []maps.Map{}
	if len(clusterIds) > 0 {
		for _, clusterId := range clusterIds {
			clusterResp, err := parentAction.RPC().NodeClusterRPC().FindEnabledNodeCluster(parentAction.AdminContext(), &pb.FindEnabledNodeClusterRequest{NodeClusterId: clusterId})
			if err != nil {
				return nil, err
			}
			var cluster = clusterResp.NodeCluster
			if cluster != nil {
				clusterMaps = append(clusterMaps, maps.Map{
					"id":   cluster.Id,
					"name": cluster.Name,
				})
			}
		}
	}
	return clusterMaps, nil
}
