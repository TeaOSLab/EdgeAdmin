package ipaddressutils

import (
	"encoding/json"
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
		var thresholdsJSON = []byte{}
		var thresholds = addr.GetSlice("thresholds")
		if len(thresholds) > 0 {
			thresholdsJSON, _ = json.Marshal(thresholds)
		}

		addrId := addr.GetInt64("id")
		if addrId > 0 {
			var isOn = false
			if !addr.Has("isOn") { // 兼容老版本
				isOn = true
			} else {
				isOn = addr.GetBool("isOn")
			}

			_, err = parentAction.RPC().NodeIPAddressRPC().UpdateNodeIPAddress(parentAction.AdminContext(), &pb.UpdateNodeIPAddressRequest{
				AddressId:      addrId,
				Ip:             addr.GetString("ip"),
				Name:           addr.GetString("name"),
				CanAccess:      addr.GetBool("canAccess"),
				IsOn:           isOn,
				ThresholdsJSON: thresholdsJSON,
			})
			if err != nil {
				return err
			}
		} else {
			_, err = parentAction.RPC().NodeIPAddressRPC().CreateNodeIPAddress(parentAction.AdminContext(), &pb.CreateNodeIPAddressRequest{
				NodeId:         nodeId,
				Role:           role,
				Name:           addr.GetString("name"),
				Ip:             addr.GetString("ip"),
				CanAccess:      addr.GetBool("canAccess"),
				ThresholdsJSON: thresholdsJSON,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
