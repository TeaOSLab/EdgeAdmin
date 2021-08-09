// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type NodeOptionsAction struct {
	actionutils.ParentAction
}

func (this *NodeOptionsAction) RunPost(params struct {
	ClusterId int64
}) {
	nodesResp, err := this.RPC().NSNodeRPC().FindAllEnabledNSNodesWithNSClusterId(this.AdminContext(), &pb.FindAllEnabledNSNodesWithNSClusterIdRequest{NsClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var nodeMaps = []maps.Map{}
	for _, node := range nodesResp.NsNodes {
		if !node.IsOn {
			continue
		}

		addressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledIPAddressesWithNodeIdRequest{
			NodeId: node.Id,
			Role:   nodeconfigs.NodeRoleDNS,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var addresses = addressesResp.Addresses
		if len(addresses) == 0 {
			continue
		}
		var addrs = []string{}
		for _, addr := range addresses {
			if addr.CanAccess {
				addrs = append(addrs, addr.Ip)
			}
		}

		nodeMaps = append(nodeMaps, maps.Map{
			"id":    node.Id,
			"name":  node.Name,
			"addrs": addrs,
		})
	}
	this.Data["nodes"] = nodeMaps

	this.Success()
}
