// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package transfer

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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type ValidateAPIAction struct {
	actionutils.ParentAction
}

func (this *ValidateAPIAction) RunPost(params struct {
	Host     string
	Port     string
	Protocol string

	Must *actions.Must
}) {
	params.Must.
		Field("newAPINodeHost", params.Host).
		Require("请输入新的API节点IP或域名").
		Field("newAPINodePort", params.Port).
		Require("请输入新的API节点端口")

	if !regexp.MustCompile(`^\d{1,5}$`).MatchString(params.Port) {
		this.FailField("newAPINodePort", "请输入正确的端口")
	}

	// 检查端口
	config, err := configs.LoadAPIConfig()
	if err != nil {
		this.Fail("加载当前平台的API配置失败：" + err.Error())
	}
	config.RPCEndpoints = []string{params.Protocol + "://" + configutils.QuoteIP(params.Host) + ":" + params.Port}
	client, err := rpc.NewRPCClient(config, false)
	if err != nil {
		this.Fail("检查API节点地址出错：" + err.Error())
	}
	defer func() {
		_ = client.Close()
	}()

	_, err = client.AdminRPC().FindAdminFullname(this.AdminContext(), &pb.FindAdminFullnameRequest{AdminId: this.AdminId()})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.Unavailable {
				this.Fail("测试新API节点失败：无法连接新的API节点：请检查：1、API节点地址和端口是否正确；2、防火墙或安全策略是否已正确设置。详细原因：" + err.Error())
			}
		}
		this.Fail("测试新API节点失败：" + err.Error())
	}

	// 所有API节点
	apiNodesResp, err := client.APINodeRPC().FindAllEnabledAPINodes(this.AdminContext(), &pb.FindAllEnabledAPINodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var apiNodes = apiNodesResp.ApiNodes
	var hosts = []string{}
	for _, node := range apiNodes {
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
	this.Data["hosts"] = hosts

	this.Success()
}
