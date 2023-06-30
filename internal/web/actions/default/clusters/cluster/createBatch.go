package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"net"
	"strconv"
	"strings"
)

type CreateBatchAction struct {
	actionutils.ParentAction
}

func (this *CreateBatchAction) Init() {
	this.Nav("", "node", "create")
	this.SecondMenu("nodes")
}

func (this *CreateBatchAction) RunGet(params struct {
	ClusterId int64
}) {
	leftMenuItems := []maps.Map{
		{
			"name":     this.Lang(codes.NodeMenu_CreateSingleNode),
			"url":      "/clusters/cluster/createNode?clusterId=" + strconv.FormatInt(params.ClusterId, 10),
			"isActive": false,
		},
		{
			"name":     this.Lang(codes.NodeMenu_CreateMultipleNodes),
			"url":      "/clusters/cluster/createBatch?clusterId=" + strconv.FormatInt(params.ClusterId, 10),
			"isActive": true,
		},
	}
	this.Data["leftMenuItems"] = leftMenuItems

	// 限额
	maxNodes, leftNodes, err := this.findNodesQuota()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["quota"] = maps.Map{
		"maxNodes":  maxNodes,
		"leftNodes": leftNodes,
	}

	this.Show()
}

func (this *CreateBatchAction) RunPost(params struct {
	ClusterId int64
	GroupId   int64
	RegionId  int64
	IpList    string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	if params.ClusterId <= 0 {
		this.Fail("请选择正确的集群")
	}

	// 校验
	// TODO 支持IP范围，比如：192.168.1.[100-105]
	realIPList := []string{}
	for _, ip := range strings.Split(params.IpList, "\n") {
		ip = strings.TrimSpace(ip)
		if len(ip) == 0 {
			continue
		}
		ip = strings.ReplaceAll(ip, " ", "")

		if net.ParseIP(ip) == nil {
			this.Fail("发现错误的IP地址：" + ip)
		}

		if lists.ContainsString(realIPList, ip) {
			continue
		}
		realIPList = append(realIPList, ip)
	}

	// 保存
	for _, ip := range realIPList {
		resp, err := this.RPC().NodeRPC().CreateNode(this.AdminContext(), &pb.CreateNodeRequest{
			Name:          ip,
			NodeClusterId: params.ClusterId,
			NodeGroupId:   params.GroupId,
			NodeRegionId:  params.RegionId,
			NodeLogin:     nil,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		nodeId := resp.NodeId
		_, err = this.RPC().NodeIPAddressRPC().CreateNodeIPAddress(this.AdminContext(), &pb.CreateNodeIPAddressRequest{
			NodeId:    nodeId,
			Role:      nodeconfigs.NodeRoleNode,
			Name:      "IP地址",
			Ip:        ip,
			CanAccess: true,
			IsUp:      true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// 创建日志
	defer this.CreateLogInfo(codes.Node_LogCreateNodeBatch)

	this.Success()
}
