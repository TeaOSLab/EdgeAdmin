package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net"
	"regexp"
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
			Require("请输入API节点端口").
			Match(`^\d+$`, "API节点端口只能是数字").
			MinLength(4, "请输入4位以上的数字端口号").
			MaxLength(5, "请输入5位以下的数字端口号")
		var newPort = types.Int(params.NewPort)
		if newPort < 1024 {
			this.FailField("newPort", "API端口号不能小于1024")
		}
		if newPort > 65534 {
			this.FailField("newPort", "API端口号不能大于65534")
		}

		if net.ParseIP(params.NewHost) == nil {
			// 检查是否为域名
			var domainReg = regexp.MustCompile(`^[\w-]+$`) // 这里暂不支持 unicode 域名
			var digitReg = regexp.MustCompile(`^\d+$`)
			var isIP = true
			for _, piece := range strings.Split(params.NewHost, ".") {
				if !domainReg.MatchString(piece) {
					this.FailField("newHost", "请输入正确的API节点主机地址")
					return
				}
				if !digitReg.MatchString(piece) {
					isIP = false
				}
			}

			if isIP {
				this.FailField("newHost", "请输入正确的API节点主机地址")
				return
			}
		}

		params.Must.
			Field("newHost", params.NewHost).
			Require("请输入API节点主机地址")

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
		RPCEndpoints: []string{params.OldProtocol + "://" + configutils.QuoteIP(params.OldHost) + ":" + params.OldPort},
		NodeId:       params.OldNodeId,
		Secret:       params.OldNodeSecret,
	}, false)
	if err != nil {
		this.FailField("oldHost", "测试API节点时出错，请检查配置，错误信息："+err.Error())
	}

	defer func() {
		_ = client.Close()
	}()

	_, err = client.APINodeRPC().FindCurrentAPINodeVersion(client.APIContext(0), &pb.FindCurrentAPINodeVersionRequest{})
	if err != nil {
		this.FailField("oldHost", "无法连接此API节点，错误信息："+err.Error())
	}

	this.Success()
}
