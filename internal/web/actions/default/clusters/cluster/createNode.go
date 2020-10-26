package cluster

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strconv"
)

// 创建节点
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
	leftMenuItems := []maps.Map{
		{
			"name":     "单个创建",
			"url":      "/clusters/cluster/createNode?clusterId=" + strconv.FormatInt(params.ClusterId, 10),
			"isActive": true,
		},
		{
			"name":     "批量创建",
			"url":      "/clusters/cluster/createBatch?clusterId=" + strconv.FormatInt(params.ClusterId, 10),
			"isActive": false,
		},
	}
	this.Data["leftMenuItems"] = leftMenuItems

	this.Show()
}

func (this *CreateNodeAction) RunPost(params struct {
	Name            string
	IpAddressesJSON []byte
	ClusterId       int64
	GrantId         int64
	SshHost         string
	SshPort         int

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称")

	if len(params.IpAddressesJSON) == 0 {
		this.Fail("请至少添加一个IP地址")
	}

	// TODO 检查cluster
	if params.ClusterId <= 0 {
		this.Fail("请选择所在集群")
	}

	// TODO 检查登录授权
	loginInfo := &pb.NodeLogin{
		Id:   0,
		Name: "SSH",
		Type: "ssh",
		Params: maps.Map{
			"grantId": params.GrantId,
			"host":    params.SshHost,
			"port":    params.SshPort,
		}.AsJSON(),
	}

	// 保存
	createResp, err := this.RPC().NodeRPC().CreateNode(this.AdminContext(), &pb.CreateNodeRequest{
		Name:      params.Name,
		ClusterId: params.ClusterId,
		Login:     loginInfo,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	nodeId := createResp.NodeId

	// IP地址
	ipAddresses := []maps.Map{}
	if len(params.IpAddressesJSON) > 0 {
		err = json.Unmarshal(params.IpAddressesJSON, &ipAddresses)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, address := range ipAddresses {
			addressId := address.GetInt64("id")
			if addressId > 0 {
				_, err = this.RPC().NodeIPAddressRPC().UpdateNodeIPAddressNodeId(this.AdminContext(), &pb.UpdateNodeIPAddressNodeIdRequest{
					AddressId: addressId,
					NodeId:    nodeId,
				})
			} else {
				_, err = this.RPC().NodeIPAddressRPC().CreateNodeIPAddress(this.AdminContext(), &pb.CreateNodeIPAddressRequest{
					NodeId:    nodeId,
					Name:      address.GetString("name"),
					Ip:        address.GetString("ip"),
					CanAccess: address.GetBool("canAccess"),
				})
			}
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	}

	this.Success()
}
