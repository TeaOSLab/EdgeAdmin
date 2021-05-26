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
		addrId := addr.GetInt64("id")
		if addrId > 0 {
			_, err = parentAction.RPC().NodeIPAddressRPC().UpdateNodeIPAddress(parentAction.AdminContext(), &pb.UpdateNodeIPAddressRequest{
				AddressId: addrId,
				Ip:        addr.GetString("ip"),
				Name:      addr.GetString("name"),
				CanAccess: addr.GetBool("canAccess"),
			})
			if err != nil {
				return err
			}
		} else {
			_, err = parentAction.RPC().NodeIPAddressRPC().CreateNodeIPAddress(parentAction.AdminContext(), &pb.CreateNodeIPAddressRequest{
				NodeId:    nodeId,
				Role:      role,
				Name:      addr.GetString("name"),
				Ip:        addr.GetString("ip"),
				CanAccess: addr.GetBool("canAccess"),
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
