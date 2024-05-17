// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package node

import (
	"encoding/json"
	"errors"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/apinodeutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"strings"
)

type UpgradePopupAction struct {
	actionutils.ParentAction
}

func (this *UpgradePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpgradePopupAction) RunGet(params struct {
	NodeId int64
}) {
	this.Data["nodeId"] = params.NodeId
	this.Data["nodeName"] = ""
	this.Data["currentVersion"] = ""
	this.Data["latestVersion"] = ""
	this.Data["result"] = ""
	this.Data["resultIsOk"] = true
	this.Data["canUpgrade"] = false
	this.Data["isUpgrading"] = false

	nodeResp, err := this.RPC().APINodeRPC().FindEnabledAPINode(this.AdminContext(), &pb.FindEnabledAPINodeRequest{ApiNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var node = nodeResp.ApiNode
	if node == nil {
		this.Data["result"] = "要升级的节点不存在"
		this.Data["resultIsOk"] = false
		this.Show()
		return
	}
	this.Data["nodeName"] = node.Name + " / [" + strings.Join(node.AccessAddrs, ", ") + "]"

	// 节点状态
	var status = &nodeconfigs.NodeStatus{}
	if len(node.StatusJSON) > 0 {
		err = json.Unmarshal(node.StatusJSON, &status)
		if err != nil {
			this.ErrorPage(errors.New("decode status failed: " + err.Error()))
			return
		}
		this.Data["currentVersion"] = status.BuildVersion
	} else {
		this.Data["result"] = "无法检测到节点当前版本"
		this.Data["resultIsOk"] = false
		this.Show()
		return
	}
	this.Data["latestVersion"] = teaconst.APINodeVersion

	if status.IsActive && len(status.BuildVersion) > 0 {
		canUpgrade, reason := apinodeutils.CanUpgrade(status.BuildVersion, status.OS, status.Arch)
		if !canUpgrade {
			this.Data["result"] = reason
			this.Data["resultIsOk"] = false
			this.Show()
			return
		}
		this.Data["canUpgrade"] = true
		this.Data["result"] = "等待升级"
		this.Data["resultIsOk"] = true
	} else {
		this.Data["result"] = "当前节点非连接状态无法远程升级"
		this.Data["resultIsOk"] = false
		this.Show()
		return
	}

	// 是否正在升级
	var oldUpgrader = apinodeutils.SharedManager.FindUpgrader(params.NodeId)
	if oldUpgrader != nil {
		this.Data["result"] = "正在升级中..."
		this.Data["resultIsOk"] = false
		this.Data["isUpgrading"] = true
	}

	this.Show()
}

func (this *UpgradePopupAction) RunPost(params struct {
	NodeId int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var manager = apinodeutils.SharedManager

	var oldUpgrader = manager.FindUpgrader(params.NodeId)
	if oldUpgrader != nil {
		this.Fail("正在升级中，无需重复提交 ...")
		return
	}

	var upgrader = apinodeutils.NewUpgrader(params.NodeId)
	manager.AddUpgrader(upgrader)
	defer func() {
		manager.RemoveUpgrader(upgrader)
	}()

	err := upgrader.Upgrade()
	if err != nil {
		this.Fail("升级失败：" + err.Error())
		return
	}

	this.Success()
}
