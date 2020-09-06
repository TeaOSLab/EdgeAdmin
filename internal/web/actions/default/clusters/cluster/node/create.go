package node

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "node", "")
	this.SecondMenu("nodes")
}

func (this *CreateAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name        string
	IPAddresses string `alias:"ipAddresses"`
	ClusterId   int64
	GrantId     int64
	SshHost     string
	SshPort     int

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称")

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
	err = json.Unmarshal([]byte(params.IPAddresses), &ipAddresses)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, address := range ipAddresses {
		addressId := address.GetInt64("id")
		_, err = this.RPC().NodeIPAddressRPC().UpdateNodeIPAddressNodeId(this.AdminContext(), &pb.UpdateNodeIPAddressNodeIdRequest{
			AddressId: addressId,
			NodeId:    nodeId,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Success()
}
