package node

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/ipAddresses/ipaddressutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "node", "update")
	this.SecondMenu("nodes")
}

func (this *UpdateAction) RunGet(params struct {
	NodeId int64
}) {
	this.Data["nodeId"] = params.NodeId

	nodeResp, err := this.RPC().NSNodeRPC().FindEnabledNSNode(this.AdminContext(), &pb.FindEnabledNSNodeRequest{NsNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	node := nodeResp.NsNode
	if node == nil {
		this.WriteString("找不到要操作的节点")
		return
	}

	var clusterMap maps.Map = nil
	if node.NsCluster != nil {
		clusterMap = maps.Map{
			"id":   node.NsCluster.Id,
			"name": node.NsCluster.Name,
		}
	}

	// IP地址
	ipAddressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledNodeIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledNodeIPAddressesWithNodeIdRequest{
		NodeId: params.NodeId,
		Role:   nodeconfigs.NodeRoleDNS,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	ipAddressMaps := []maps.Map{}
	for _, addr := range ipAddressesResp.NodeIPAddresses {
		ipAddressMaps = append(ipAddressMaps, maps.Map{
			"id":        addr.Id,
			"name":      addr.Name,
			"ip":        addr.Ip,
			"canAccess": addr.CanAccess,
			"isOn":      addr.IsOn,
			"isUp":      addr.IsUp,
		})
	}

	// 登录信息
	var loginMap maps.Map = nil
	if node.NodeLogin != nil {
		loginParams := maps.Map{}
		if len(node.NodeLogin.Params) > 0 {
			err = json.Unmarshal(node.NodeLogin.Params, &loginParams)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		grantMap := maps.Map{}
		grantId := loginParams.GetInt64("grantId")
		if grantId > 0 {
			grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: grantId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if grantResp.NodeGrant != nil {
				grantMap = maps.Map{
					"id":         grantResp.NodeGrant.Id,
					"name":       grantResp.NodeGrant.Name,
					"method":     grantResp.NodeGrant.Method,
					"methodName": grantutils.FindGrantMethodName(grantResp.NodeGrant.Method),
					"username":   grantResp.NodeGrant.Username,
				}
			}
		}

		loginMap = maps.Map{
			"id":     node.NodeLogin.Id,
			"name":   node.NodeLogin.Name,
			"type":   node.NodeLogin.Type,
			"params": loginParams,
			"grant":  grantMap,
		}
	}

	this.Data["node"] = maps.Map{
		"id":          node.Id,
		"name":        node.Name,
		"ipAddresses": ipAddressMaps,
		"cluster":     clusterMap,
		"isOn":        node.IsOn,
		"login":       loginMap,
	}

	// 所有集群
	resp, err := this.RPC().NSClusterRPC().FindAllEnabledNSClusters(this.AdminContext(), &pb.FindAllEnabledNSClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range resp.NsClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	LoginId int64
	GrantId int64
	SshHost string
	SshPort int

	NodeId          int64
	Name            string
	IPAddressesJSON []byte `alias:"ipAddressesJSON"`
	ClusterId       int64
	IsOn            bool

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改节点 %d", params.NodeId)

	if params.NodeId <= 0 {
		this.Fail("要操作的节点不存在")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称")

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
	if len(params.IPAddressesJSON) > 0 {
		err := json.Unmarshal(params.IPAddressesJSON, &ipAddresses)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	if len(ipAddresses) == 0 {
		this.Fail("请至少输入一个IP地址")
	}

	// TODO 检查登录授权
	loginInfo := &pb.NodeLogin{
		Id:   params.LoginId,
		Name: "SSH",
		Type: "ssh",
		Params: maps.Map{
			"grantId": params.GrantId,
			"host":    params.SshHost,
			"port":    params.SshPort,
		}.AsJSON(),
	}

	// 保存
	_, err = this.RPC().NSNodeRPC().UpdateNSNode(this.AdminContext(), &pb.UpdateNSNodeRequest{
		NsNodeId:    params.NodeId,
		Name:        params.Name,
		NsClusterId: params.ClusterId,
		IsOn:        params.IsOn,
		NodeLogin:   loginInfo,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 禁用老的IP地址
	_, err = this.RPC().NodeIPAddressRPC().DisableAllNodeIPAddressesWithNodeId(this.AdminContext(), &pb.DisableAllNodeIPAddressesWithNodeIdRequest{
		NodeId: params.NodeId,
		Role:   nodeconfigs.NodeRoleDNS,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 添加新的IP地址
	err = ipaddressutils.UpdateNodeIPAddresses(this.Parent(), params.NodeId, nodeconfigs.NodeRoleDNS, params.IPAddressesJSON)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
