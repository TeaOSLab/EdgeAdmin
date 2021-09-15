package cluster

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

// CreateNodeAction 创建节点
type CreateNodeAction struct {
	actionutils.ParentAction
}

func (this *CreateNodeAction) Init() {
	this.Nav("", "node", "create")
	this.SecondMenu("nodes")
}

func (this *CreateNodeAction) RunGet(params struct {
	ClusterId int64
}) {
	this.Show()
}

func (this *CreateNodeAction) RunPost(params struct {
	Name            string
	IpAddressesJSON []byte
	ClusterId       int64

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称")

	if len(params.IpAddressesJSON) == 0 {
		this.Fail("请至少添加一个IP地址")
	}

	// 检查cluster
	if params.ClusterId <= 0 {
		this.Fail("请选择所在集群")
	}
	clusterResp, err := this.RPC().NSClusterRPC().FindEnabledNSCluster(this.AdminContext(), &pb.FindEnabledNSClusterRequest{NsClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if clusterResp.NsCluster == nil {
		this.Fail("选择的集群不存在")
	}

	// IP地址
	ipAddresses := []maps.Map{}
	if len(params.IpAddressesJSON) > 0 {
		err := json.Unmarshal(params.IpAddressesJSON, &ipAddresses)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	if len(ipAddresses) == 0 {
		this.Fail("请至少输入一个IP地址")
	}

	// 保存
	createResp, err := this.RPC().NSNodeRPC().CreateNSNode(this.AdminContext(), &pb.CreateNSNodeRequest{
		Name:          params.Name,
		NodeClusterId: params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	nodeId := createResp.NsNodeId

	// IP地址
	for _, addrMap := range ipAddresses {
		addressId := addrMap.GetInt64("id")
		if addressId > 0 {
			_, err = this.RPC().NodeIPAddressRPC().UpdateNodeIPAddressNodeId(this.AdminContext(), &pb.UpdateNodeIPAddressNodeIdRequest{
				NodeIPAddressId: addressId,
				NodeId:          nodeId,
			})
		} else {
			_, err = this.RPC().NodeIPAddressRPC().CreateNodeIPAddress(this.AdminContext(), &pb.CreateNodeIPAddressRequest{
				NodeId:    nodeId,
				Role:      nodeconfigs.NodeRoleDNS,
				Name:      addrMap.GetString("name"),
				Ip:        addrMap.GetString("ip"),
				CanAccess: addrMap.GetBool("canAccess"),
				IsUp:      addrMap.GetBool("isUp"),
			})
		}
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "创建域名服务节点 %d", nodeId)

	this.Success()
}
