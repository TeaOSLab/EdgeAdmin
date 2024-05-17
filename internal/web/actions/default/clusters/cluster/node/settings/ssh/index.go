// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package ssh

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net"
	"regexp"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "update")
	this.SecondMenu("ssh")
}

func (this *IndexAction) RunGet(params struct {
	NodeId int64
}) {
	node, err := nodeutils.InitNodeInfo(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["hostIsAutoFilled"] = false

	// 登录信息
	var loginMap maps.Map = nil
	if node.NodeLogin != nil {
		var loginParams = maps.Map{}
		if len(node.NodeLogin.Params) > 0 {
			err = json.Unmarshal(node.NodeLogin.Params, &loginParams)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		var grantMap = maps.Map{}
		var grantId = loginParams.GetInt64("grantId")
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
					"methodName": grantutils.FindGrantMethodName(grantResp.NodeGrant.Method, this.LangCode()),
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

	if loginMap == nil {
		addressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledNodeIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledNodeIPAddressesWithNodeIdRequest{NodeId: node.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(addressesResp.NodeIPAddresses) > 0 {
			this.Data["hostIsAutoFilled"] = true
			loginMap = maps.Map{
				"id":   0,
				"name": "",
				"type": "ssh",
				"params": maps.Map{
					"host":    addressesResp.NodeIPAddresses[0].Ip,
					"port":    22,
					"grantId": 0,
				},
				"grant": nil,
			}
		}
	} else {
		var loginParams = loginMap.GetMap("params")
		if len(loginParams.GetString("host")) == 0 {
			addressesResp, err := this.RPC().NodeIPAddressRPC().FindAllEnabledNodeIPAddressesWithNodeId(this.AdminContext(), &pb.FindAllEnabledNodeIPAddressesWithNodeIdRequest{NodeId: node.Id})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if len(addressesResp.NodeIPAddresses) > 0 {
				this.Data["hostIsAutoFilled"] = true
				loginParams["host"] = addressesResp.NodeIPAddresses[0].Ip
			}
		}

		if loginParams.GetInt("port") == 0 {
			loginParams["port"] = 22
		}

		loginMap["params"] = loginParams
	}

	var nodeMap = this.Data["node"].(maps.Map)
	nodeMap["login"] = loginMap

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	NodeId int64

	LoginId int64
	GrantId int64
	SshHost string
	SshPort int

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.NodeSSH_LogUpdateNodeSSH, params.NodeId)

	// 检查IP地址
	if regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`).MatchString(params.SshHost) && net.ParseIP(params.SshHost) == nil {
		this.Fail("SSH主机地址 '" + params.SshHost + "' IP格式错误")
	}

	// TODO 检查登录授权
	var loginInfo = &pb.NodeLogin{
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
		NodeLogin: loginInfo,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
