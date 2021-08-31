// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package addr

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/ip-addrs/ipaddrutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "addr")
}

func (this *IndexAction) RunGet(params struct {
	AddrId int64
}) {
	addr, err := ipaddrutils.InitIPAddr(this.Parent(), params.AddrId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

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
		this.ErrorPage(errors.New("node or cluster is not available"))
		return
	}

	this.Data["addr"] = maps.Map{
		"id":          addr.Id,
		"name":        addr.Name,
		"description": addr.Description,
		"ip":          addr.Ip,
		"canAccess":   addr.CanAccess,
		"isOn":        addr.IsOn,
		"isUp":        addr.IsUp,
		"thresholds":  thresholds,
		"node":        maps.Map{"id": node.Id, "name": node.Name},
		"cluster":     maps.Map{"id": node.NodeCluster.Id, "name": node.NodeCluster.Name},
	}

	this.Show()
}
