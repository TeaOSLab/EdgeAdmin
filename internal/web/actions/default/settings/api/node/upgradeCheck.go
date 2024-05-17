// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

// UpgradeCheckAction 检查升级结果
type UpgradeCheckAction struct {
	actionutils.ParentAction
}

func (this *UpgradeCheckAction) Init() {
	this.Nav("", "", "")
}

func (this *UpgradeCheckAction) RunPost(params struct {
	NodeId int64
}) {
	this.Data["isOk"] = false

	nodeResp, err := this.RPC().APINodeRPC().FindEnabledAPINode(this.AdminContext(), &pb.FindEnabledAPINodeRequest{ApiNodeId: params.NodeId})
	if err != nil {
		this.Success()
		return
	}

	var node = nodeResp.ApiNode
	if node == nil || len(node.AccessAddrs) == 0 {
		this.Success()
		return
	}

	apiConfig, err := configs.LoadAPIConfig()
	if err != nil {
		this.Success()
		return
	}

	var newAPIConfig = apiConfig.Clone()
	newAPIConfig.RPCEndpoints = node.AccessAddrs
	rpcClient, err := rpc.NewRPCClient(newAPIConfig, false)
	if err != nil {
		this.Success()
		return
	}

	versionResp, err := rpcClient.APINodeRPC().FindCurrentAPINodeVersion(rpcClient.Context(0), &pb.FindCurrentAPINodeVersionRequest{})
	if err != nil {
		this.Success()
		return
	}

	if versionResp.Version != teaconst.Version {
		this.Success()
		return
	}

	this.Data["isOk"] = true

	this.Success()
}
