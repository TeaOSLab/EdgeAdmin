// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package confirm

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"net/url"
	"strings"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	var endpoints = []string{}

	config, err := configs.LoadAPIConfig()
	if err == nil {
		endpoints = config.RPCEndpoints
		this.Data["nodeId"] = config.NodeId
		this.Data["secret"] = config.Secret
		this.Data["canInstall"] = false
	} else {
		this.Data["nodeId"] = ""
		this.Data["secret"] = ""
		this.Data["canInstall"] = true
	}

	if len(endpoints) == 0 {
		endpoints = []string{""} // 初始化一个空的
	}

	this.Data["endpoints"] = endpoints

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	Endpoints []string
	NodeId    string
	Secret    string

	Must *actions.Must
}) {
	var endpoints = []string{}
	for _, endpoint := range params.Endpoints {
		if len(endpoint) > 0 {
			u, err := url.Parse(endpoint)
			if err != nil {
				this.Fail("API节点地址'" + endpoint + "'格式错误")
			}
			endpoint = u.Scheme + "://" + u.Host
			if u.Scheme != "http" && u.Scheme != "https" {
				this.Fail("API节点地址'" + endpoint + "'中的协议错误，目前只支持http或者https")
			}
			switch u.Scheme {
			case "http":
				if len(u.Port()) == 0 {
					endpoint += ":80"
				}
			case "https":
				if len(u.Port()) == 0 {
					endpoint += ":443"
				}
			}

			// 检测是否连接
			var config = &configs.APIConfig{}
			config.NodeId = params.NodeId
			config.Secret = params.Secret
			config.RPCEndpoints = []string{endpoint}
			client, err := rpc.NewRPCClient(config, false)
			if err != nil {
				this.Fail("尝试配置RPC发生错误：" + err.Error())
				return
			}
			resp, err := client.APINodeRPC().FindCurrentAPINodeVersion(client.Context(0), &pb.FindCurrentAPINodeVersionRequest{})
			if err != nil {
				_ = client.Close()

				if strings.Contains(err.Error(), "wrong token role") {
					this.Fail("你输入的NodeId和Secret为其他节点的配置信息，不是管理系统的配置信息，所以无法使用；请从管理系统的配置目录下找到管理系统的配置信息并填入。如果你不知道如何查找，请刷新当前页面，使用默认填写的NodeId和Secret提交。")
				} else {
					this.Fail("无法连接你填入的API节点地址，请检查协议、IP和端口是否正确，错误信息：" + err.Error())
				}
				return
			}

			if resp != nil && resp.Role != "admin" {
				this.Fail("你输入的NodeId和Secret为API节点的配置信息，不是管理系统的配置信息，所以无法使用；请从管理系统的配置目录下找到管理系统的配置信息并填入")
				return
			}
			_ = client.Close()

			endpoints = append(endpoints, endpoint)
		}
	}

	if len(endpoints) == 0 {
		this.Fail("请输入至少一个API节点地址")
	}

	if len(params.NodeId) == 0 {
		this.Fail("请输入NodeId")
	}
	if len(params.Secret) == 0 {
		this.Fail("请输入Secret")
	}

	// 创建配置文件
	config, err := configs.LoadAPIConfig()
	if err != nil {
		config = &configs.APIConfig{}
	}
	config.NodeId = params.NodeId
	config.Secret = params.Secret
	config.RPCEndpoints = endpoints
	config.RPCDisableUpdate = true
	err = config.WriteFile(Tea.ConfigFile(configs.ConfigFileName))
	if err != nil {
		this.Fail("配置保存失败：" + err.Error())
	}

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		this.Fail("RPC配置无法读取：" + err.Error())
	}
	err = rpcClient.UpdateConfig(config)
	if err != nil {
		this.Fail("重载RPC配置失败：" + err.Error())
	}

	this.Success()
}
