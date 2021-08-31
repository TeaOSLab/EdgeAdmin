// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipaddrs

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "index")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
	UpState   int8
	Keyword   string
}) {
	this.Data["clusterId"] = params.ClusterId
	this.Data["upState"] = params.UpState
	this.Data["keyword"] = params.Keyword

	countResp, err := this.RPC().NodeIPAddressRPC().CountAllEnabledIPAddresses(this.AdminContext(), &pb.CountAllEnabledIPAddressesRequest{
		NodeClusterId: params.ClusterId,
		Role:          nodeconfigs.NodeRoleNode,
		UpState:       int32(params.UpState),
		Keyword:       params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var count = countResp.Count
	var page = this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	addrsResp, err := this.RPC().NodeIPAddressRPC().ListEnabledIPAddresses(this.AdminContext(), &pb.ListEnabledIPAddressesRequest{
		NodeClusterId: params.ClusterId,
		Role:          nodeconfigs.NodeRoleNode,
		UpState:       int32(params.UpState),
		Keyword:       params.Keyword,
		Offset:        page.Offset,
		Size:          page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var addrMaps = []maps.Map{}
	for _, addr := range addrsResp.NodeIPAddresses {
		var thresholds = []*nodeconfigs.NodeValueThresholdConfig{}
		if len(addr.ThresholdsJSON) > 0 {
			err = json.Unmarshal(addr.ThresholdsJSON, &thresholds)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		nodeResp, err := this.RPC().NodeRPC().FindEnabledBasicNode(this.AdminContext(), &pb.FindEnabledBasicNodeRequest{NodeId: addr.NodeId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var node = nodeResp.Node
		if node == nil || node.NodeCluster == nil {
			continue
		}

		addrMaps = append(addrMaps, maps.Map{
			"id":            addr.Id,
			"name":          addr.Name,
			"description":   addr.Description,
			"ip":            addr.Ip,
			"canAccess":     addr.CanAccess,
			"isOn":          addr.IsOn,
			"isUp":          addr.IsUp,
			"hasThresholds": len(thresholds) > 0,
			"node":          maps.Map{"id": node.Id, "name": node.Name},
			"cluster":       maps.Map{"id": node.NodeCluster.Id, "name": node.NodeCluster.Name},
		})
	}
	this.Data["addrs"] = addrMaps

	this.Show()
}
