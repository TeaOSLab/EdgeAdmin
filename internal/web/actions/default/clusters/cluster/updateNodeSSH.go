package cluster

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net"
	"regexp"
)

type UpdateNodeSSHAction struct {
	actionutils.ParentAction
}

func (this *UpdateNodeSSHAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateNodeSSHAction) RunGet(params struct {
	NodeId int64
}) {
	nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if nodeResp.Node == nil {
		this.NotFound("node", params.NodeId)
		return
	}

	var node = nodeResp.Node
	this.Data["node"] = maps.Map{
		"id":   node.Id,
		"name": node.Name,
	}
	if nodeResp.Node.NodeCluster != nil {
		this.Data["clusterId"] = nodeResp.Node.NodeCluster.Id
	} else {
		this.Data["clusterId"] = 0
	}

	// SSH
	var loginParams = maps.Map{
		"host":    "",
		"port":    "",
		"grantId": 0,
	}
	this.Data["loginId"] = 0
	if node.NodeLogin != nil {
		this.Data["loginId"] = node.NodeLogin.Id
		if len(node.NodeLogin.Params) > 0 {
			err = json.Unmarshal(node.NodeLogin.Params, &loginParams)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	}

	if len(loginParams.GetString("host")) == 0 {
		addressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledNodeIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledNodeIPAddressesWithNodeIdRequest{NodeId: node.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(addressesResp.NodeIPAddresses) > 0 {
			loginParams["host"] = addressesResp.NodeIPAddresses[0].Ip
		}
	}

	this.Data["params"] = loginParams

	// 认证信息
	var grantId = loginParams.GetInt64("grantId")
	grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: grantId})
	if err != nil {
		this.ErrorPage(err)
	}
	var grantMap maps.Map = nil
	if grantResp.NodeGrant != nil {
		grantMap = maps.Map{
			"id":         grantResp.NodeGrant.Id,
			"name":       grantResp.NodeGrant.Name,
			"method":     grantResp.NodeGrant.Method,
			"methodName": grantutils.FindGrantMethodName(grantResp.NodeGrant.Method, this.LangCode()),
		}
	}
	this.Data["grant"] = grantMap

	this.Show()
}

func (this *UpdateNodeSSHAction) RunPost(params struct {
	NodeId  int64
	LoginId int64
	SshHost string
	SshPort int
	GrantId int64

	Must *actions.Must
}) {
	params.Must.
		Field("sshHost", params.SshHost).
		Require("请输入SSH主机地址").
		Field("sshPort", params.SshPort).
		Gt(0, "SSH主机端口需要大于0").
		Lt(65535, "SSH主机端口需要小于65535")

	if params.GrantId <= 0 {
		this.Fail("需要选择或填写至少一个认证信息")
	}

	// 检查IP地址
	if regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`).MatchString(params.SshHost) && net.ParseIP(params.SshHost) == nil {
		this.Fail("SSH主机地址 '" + params.SshHost + "' IP格式错误")
	}

	var login = &pb.NodeLogin{
		Id:   params.LoginId,
		Name: "SSH",
		Type: "ssh",
		Params: maps.Map{
			"grantId": params.GrantId,
			"host":    params.SshHost,
			"port":    params.SshPort,
		}.AsJSON(),
	}

	_, err := this.RPC().NodeRPC().UpdateNodeLogin(this.AdminContext(), &pb.UpdateNodeLoginRequest{
		NodeId:    params.NodeId,
		NodeLogin: login,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLogInfo(codes.NodeSSH_LogUpdateNodeSSH, params.NodeId)

	this.Success()
}
