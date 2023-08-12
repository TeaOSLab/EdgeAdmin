package recovers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type ValidateApiAction struct {
	actionutils.ParentAction
}

func (this *ValidateApiAction) RunPost(params struct {
	Protocol   string
	Host       string
	Port       string
	NodeId     string
	NodeSecret string

	Must *actions.Must
}) {
	params.NodeId = strings.Trim(params.NodeId, "\"' ")
	params.NodeSecret = strings.Trim(params.NodeSecret, "\"' ")

	// 使用已有的API节点
	params.Must.
		Field("host", params.Host).
		Require("请输入主机地址").
		Field("port", params.Port).
		Require("请输入服务端口").
		Match(`^\d+$`, "服务端口只能是数字").
		Field("nodeId", params.NodeId).
		Require("请输入节点nodeId").
		Field("nodeSecret", params.NodeSecret).
		Require("请输入节点secret")
	client, err := rpc.NewRPCClient(&configs.APIConfig{
		RPCEndpoints: []string{params.Protocol + "://" + configutils.QuoteIP(params.Host) + ":" + params.Port},
		NodeId:       params.NodeId,
		Secret:       params.NodeSecret,
	}, false)
	if err != nil {
		this.FailField("host", "测试API节点时出错，请检查配置，错误信息："+err.Error())
	}

	defer func() {
		_ = client.Close()
	}()

	_, err = client.APINodeRPC().FindCurrentAPINodeVersion(client.APIContext(0), &pb.FindCurrentAPINodeVersionRequest{})
	if err != nil {
		this.FailField("host", "无法连接此API节点，错误信息："+err.Error())
	}

	// API节点列表
	nodesResp, err := client.APINodeRPC().FindAllEnabledAPINodes(client.Context(0), &pb.FindAllEnabledAPINodesRequest{})
	if err != nil {
		this.Fail("获取API节点列表失败，错误信息：" + err.Error())
	}
	var hosts = []string{}
	for _, node := range nodesResp.ApiNodes {
		if !node.IsOn {
			continue
		}

		// http
		if len(node.HttpJSON) > 0 {
			var config = &serverconfigs.HTTPProtocolConfig{}
			err = json.Unmarshal(node.HttpJSON, config)
			if err != nil {
				this.Fail("读取节点HTTP信息失败：" + err.Error())
			}
			for _, listen := range config.Listen {
				if len(listen.Host) > 0 && !lists.ContainsString(hosts, listen.Host) {
					hosts = append(hosts, listen.Host)
				}
			}
		}

		// https
		if len(node.HttpsJSON) > 0 {
			var config = &serverconfigs.HTTPSProtocolConfig{}
			err = json.Unmarshal(node.HttpsJSON, config)
			if err != nil {
				this.Fail("读取节点HTTPS信息失败：" + err.Error())
			}
			for _, listen := range config.Listen {
				if len(listen.Host) > 0 && !lists.ContainsString(hosts, listen.Host) {
					hosts = append(hosts, listen.Host)
				}
			}
		}

		// restHTTP
		if len(node.RestHTTPJSON) > 0 {
			var config = &serverconfigs.HTTPProtocolConfig{}
			err = json.Unmarshal(node.RestHTTPJSON, config)
			if err != nil {
				this.Fail("读取节点REST HTTP信息失败：" + err.Error())
			}
			for _, listen := range config.Listen {
				if len(listen.Host) > 0 && !lists.ContainsString(hosts, listen.Host) {
					hosts = append(hosts, listen.Host)
				}
			}
		}

		// restHTTPS
		if len(node.RestHTTPSJSON) > 0 {
			var config = &serverconfigs.HTTPSProtocolConfig{}
			err = json.Unmarshal(node.RestHTTPSJSON, config)
			if err != nil {
				this.Fail("读取节点REST HTTPS信息失败：" + err.Error())
			}
			for _, listen := range config.Listen {
				if len(listen.Host) > 0 && !lists.ContainsString(hosts, listen.Host) {
					hosts = append(hosts, listen.Host)
				}
			}
		}

		// access addrs
		if len(node.AccessAddrsJSON) > 0 {
			var addrs []*serverconfigs.NetworkAddressConfig
			err = json.Unmarshal(node.AccessAddrsJSON, &addrs)
			if err != nil {
				this.Fail("读取节点访问地址失败：" + err.Error())
			}
			for _, addr := range addrs {
				if len(addr.Host) > 0 && !lists.ContainsString(hosts, addr.Host) {
					hosts = append(hosts, addr.Host)
				}
			}
		}
	}

	this.Data["apiNode"] = maps.Map{
		"protocol":   params.Protocol,
		"host":       params.Host,
		"port":       params.Port,
		"nodeId":     params.NodeId,
		"nodeSecret": params.NodeSecret,
		"hosts":      hosts,
	}

	this.Success()
}
