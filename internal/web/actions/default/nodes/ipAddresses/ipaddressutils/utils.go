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
	addresses := []maps.Map{}
	err := json.Unmarshal(ipAddressesJSON, &addresses)
	if err != nil {
		return err
	}
	for _, addr := range addresses {
		var resultAddrIds = []int64{}
		addrId := addr.GetInt64("id")
		if addrId > 0 {
			resultAddrIds = append(resultAddrIds, addrId)

			var isOn = false
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
					NodeId:    nodeId,
					Role:      role,
					Name:      addr.GetString("name"),
					Ip:        result[0],
					CanAccess: addr.GetBool("canAccess"),
					IsUp:      addr.GetBool("isUp"),
				})
				if err != nil {
					return err
				}
				addrId = createResp.NodeIPAddressId
				resultAddrIds = append(resultAddrIds, addrId)
			} else if len(result) > 1 {
				// 批量创建
				createResp, err := parentAction.RPC().NodeIPAddressRPC().CreateNodeIPAddresses(parentAction.AdminContext(), &pb.CreateNodeIPAddressesRequest{
					NodeId:     nodeId,
					Role:       role,
					Name:       addr.GetString("name"),
					IpList:     result,
					CanAccess:  addr.GetBool("canAccess"),
					IsUp:       addr.GetBool("isUp"),
					GroupValue: ipStrings,
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
