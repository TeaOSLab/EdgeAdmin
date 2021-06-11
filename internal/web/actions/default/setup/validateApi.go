package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net"
	"strings"
)

type ValidateApiAction struct {
	actionutils.ParentAction
}

func (this *ValidateApiAction) RunPost(params struct {
	Mode          string
	NewPort       string
	NewHost       string
	OldProtocol   string
	OldHost       string
	OldPort       string
	OldNodeId     string
	OldNodeSecret string

	Must *actions.Must
}) {
	params.OldNodeId = strings.Trim(params.OldNodeId, "\"' ")
	params.OldNodeSecret = strings.Trim(params.OldNodeSecret, "\"' ")

	this.Data["apiNode"] = maps.Map{
		"mode": params.Mode,

		"newPort": params.NewPort,
		"newHost": params.NewHost,

		"oldProtocol":   params.OldProtocol,
		"oldHost":       params.OldHost,
		"oldPort":       params.OldPort,
		"oldNodeId":     params.OldNodeId,
		"oldNodeSecret": params.OldNodeSecret,
	}

	if params.Mode == "new" {
		params.Must.
			Field("newPort", params.NewPort).
			Require("请输入节点端口").
			Match(`^\d+$`, "节点端口只能是数字").
			MinLength(4, "请输入4位以上的数字").
			MaxLength(5, "请输入5位以下的数字")
		newPort := types.Int(params.NewPort)
		if newPort < 1024 {
			this.FailField("newPort", "端口号不能小于1024")
		}
		if newPort > 65534 {
			this.FailField("newPort", "端口号不能大于65534")
		}

		if net.ParseIP(params.NewHost) == nil {
			this.FailField("newHost", "请输入正确的节点主机地址")
		}

		listener, err := net.Listen("tcp", params.NewHost+":"+params.NewPort)
		if err != nil {
			if strings.Contains(err.Error(), "in use") {
				this.FailField("newPort", "端口 \""+params.NewPort+"\" 已经被别的服务占用，请关闭正在使用此端口的其他应用程序，或者使用另外一个端口")
			} else {
				this.FailField("newHost", "无法正确绑定端口地址，请检查主机地址："+err.Error())
			}
		}
		_ = listener.Close()

		params.Must.
			Field("newHost", params.NewHost).
			Require("请输入节点主机地址")

		this.Success()
		return
	}

	// 使用已有的API节点
	params.Must.
		Field("oldHost", params.OldHost).
		Require("请输入主机地址").
		Field("oldPort", params.OldPort).
		Require("请输入服务端口").
		Match(`^\d+$`, "服务端口只能是数字").
		Field("oldNodeId", params.OldNodeId).
		Require("请输入节点nodeId").
		Field("oldNodeSecret", params.OldNodeSecret).
		Require("请输入节点secret")
	client, err := rpc.NewRPCClient(&configs.APIConfig{
		RPC: struct {
			Endpoints []string `yaml:"endpoints"`
		}{
			Endpoints: []string{params.OldProtocol + "://" + params.OldHost + ":" + params.OldPort},
		},
		NodeId: params.OldNodeId,
		Secret: params.OldNodeSecret,
	})
	if err != nil {
		this.FailField("oldHost", "测试API节点时出错，请检查配置，错误信息："+err.Error())
	}
	_, err = client.APINodeRPC().FindCurrentAPINodeVersion(client.APIContext(0), &pb.FindCurrentAPINodeVersionRequest{})
	if err != nil {
		this.FailField("oldHost", "无法连接此API节点，错误信息："+err.Error())
	}

	this.Success()
}
