// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package system

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "update")
	this.SecondMenu("system")
}

func (this *IndexAction) RunGet(params struct {
	NodeId int64
}) {
	node, err := nodeutils.InitNodeInfo(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 获取节点信息
	var nodeMap = this.Data["node"].(maps.Map)
	nodeMap["maxCPU"] = node.MaxCPU

	// DNS
	dnsResolverResp, err := this.RPC().NodeRPC().FindNodeDNSResolver(this.AdminContext(), &pb.FindNodeDNSResolverRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var dnsResolverConfig = nodeconfigs.DefaultDNSResolverConfig()
	if len(dnsResolverResp.DnsResolverJSON) > 0 {
		err = json.Unmarshal(dnsResolverResp.DnsResolverJSON, dnsResolverConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["dnsResolverConfig"] = dnsResolverConfig

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	NodeId int64
	MaxCPU int32

	DnsResolverJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改节点 %d 系统信息", params.NodeId)

	if params.MaxCPU < 0 {
		this.Fail("CPU线程数不能小于0")
	}

	_, err := this.RPC().NodeRPC().UpdateNodeSystem(this.AdminContext(), &pb.UpdateNodeSystemRequest{
		NodeId: params.NodeId,
		MaxCPU: params.MaxCPU,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var dnsResolverConfig = nodeconfigs.DefaultDNSResolverConfig()
	err = json.Unmarshal(params.DnsResolverJSON, dnsResolverConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	err = dnsResolverConfig.Init()
	if err != nil {
		this.Fail("校验DNS解析配置失败：" + err.Error())
	}

	_, err = this.RPC().NodeRPC().UpdateNodeDNSResolver(this.AdminContext(), &pb.UpdateNodeDNSResolverRequest{
		NodeId:          params.NodeId,
		DnsResolverJSON: params.DnsResolverJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
