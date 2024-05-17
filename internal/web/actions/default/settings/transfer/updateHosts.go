// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package transfer

import (
	"bytes"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/lists"
)

type UpdateHostsAction struct {
	actionutils.ParentAction
}

func (this *UpdateHostsAction) RunPost(params struct {
	Protocol string
	Host     string
	Port     string

	OldHosts []string
	NewHosts []string
}) {
	if len(params.OldHosts) != len(params.NewHosts) {
		this.Fail("参数配置错误，请刷新页面后重试")
	}

	// 检查端口
	config, err := configs.LoadAPIConfig()
	if err != nil {
		this.Fail("加载当前平台的API配置失败：" + err.Error())
	}
	var apiURL = params.Protocol + "://" + configutils.QuoteIP(params.Host) + ":" + params.Port
	config.RPCEndpoints = []string{apiURL}
	client, err := rpc.NewRPCClient(config, false)
	if err != nil {
		this.Fail("检查API节点地址出错：" + err.Error())
	}
	defer func() {
		_ = client.Close()
	}()

	if err != nil {
		this.FailField("host", "测试API节点时出错，请检查配置，错误信息："+err.Error())
	}
	_, err = client.APINodeRPC().FindCurrentAPINodeVersion(client.Context(0), &pb.FindCurrentAPINodeVersionRequest{})
	if err != nil {
		this.FailField("host", "无法连接此API节点，错误信息："+err.Error())
	}

	defer func() {
		_ = client.Close()
	}()

	// API节点列表
	nodesResp, err := client.APINodeRPC().FindAllEnabledAPINodes(client.Context(0), &pb.FindAllEnabledAPINodesRequest{})
	if err != nil {
		this.Fail("获取API节点列表失败，错误信息：" + err.Error())
	}
	var endpoints = []string{}
	for _, node := range nodesResp.ApiNodes {
		if !node.IsOn {
			continue
		}

		// http
		if len(node.HttpJSON) > 0 {
			for index, oldHost := range params.OldHosts {
				if len(params.NewHosts[index]) == 0 {
					continue
				}
				node.HttpJSON = bytes.ReplaceAll(node.HttpJSON, []byte("\""+oldHost+"\""), []byte("\""+params.NewHosts[index]+"\""))
			}
		}

		// https
		if len(node.HttpsJSON) > 0 {
			for index, oldHost := range params.OldHosts {
				if len(params.NewHosts[index]) == 0 {
					continue
				}
				node.HttpsJSON = bytes.ReplaceAll(node.HttpsJSON, []byte("\""+oldHost+"\""), []byte("\""+params.NewHosts[index]+"\""))
			}
		}

		// restHTTP
		if len(node.RestHTTPJSON) > 0 {
			for index, oldHost := range params.OldHosts {
				if len(params.NewHosts[index]) == 0 {
					continue
				}
				node.RestHTTPJSON = bytes.ReplaceAll(node.RestHTTPJSON, []byte("\""+oldHost+"\""), []byte("\""+params.NewHosts[index]+"\""))
			}
		}

		// restHTTPS
		if len(node.RestHTTPSJSON) > 0 {
			for index, oldHost := range params.OldHosts {
				if len(params.NewHosts[index]) == 0 {
					continue
				}
				node.RestHTTPSJSON = bytes.ReplaceAll(node.RestHTTPSJSON, []byte("\""+oldHost+"\""), []byte("\""+params.NewHosts[index]+"\""))
			}
		}

		// access addrs
		if len(node.AccessAddrsJSON) > 0 {
			for index, oldHost := range params.OldHosts {
				if len(params.NewHosts[index]) == 0 {
					continue
				}
				node.AccessAddrsJSON = bytes.ReplaceAll(node.AccessAddrsJSON, []byte("\""+oldHost+"\""), []byte("\""+params.NewHosts[index]+"\""))
			}

			var addrs []*serverconfigs.NetworkAddressConfig
			err = json.Unmarshal(node.AccessAddrsJSON, &addrs)
			if err != nil {
				this.Fail("读取节点访问地址失败：" + err.Error())
			}
			for _, addr := range addrs {
				err = addr.Init()
				if err != nil {
					// 暂时不提示错误
					continue
				}
				for _, a := range addr.FullAddresses() {
					if !lists.ContainsString(endpoints, a) {
						endpoints = append(endpoints, a)
					}
				}
			}
		}

		// 保存
		_, err = client.APINodeRPC().UpdateAPINode(client.Context(0), &pb.UpdateAPINodeRequest{
			ApiNodeId:       node.Id,
			Name:            node.Name,
			Description:     node.Description,
			HttpJSON:        node.HttpJSON,
			HttpsJSON:       node.HttpsJSON,
			AccessAddrsJSON: node.AccessAddrsJSON,
			IsOn:            node.IsOn,
			RestIsOn:        node.RestIsOn,
			RestHTTPJSON:    node.RestHTTPJSON,
			RestHTTPSJSON:   node.RestHTTPSJSON,
		})
		if err != nil {
			this.Fail("保存API节点信息失败：" + err.Error())
		}
	}

	this.Success()
}
