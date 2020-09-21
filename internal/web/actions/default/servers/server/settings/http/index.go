package http

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("http")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	httpConfig := &serverconfigs.HTTPProtocolConfig{}
	if len(server.HttpJSON) > 0 {
		err := json.Unmarshal(server.HttpJSON, httpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		httpConfig.IsOn = true
	}

	this.Data["serverType"] = server.Type
	this.Data["httpConfig"] = maps.Map{
		"isOn":      httpConfig.IsOn,
		"addresses": httpConfig.Listen,
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId  int64
	Addresses string

	Must *actions.Must
}) {
	addresses := []*serverconfigs.NetworkAddressConfig{}
	err := json.Unmarshal([]byte(params.Addresses), &addresses)
	if err != nil {
		this.Fail("端口地址解析失败：" + err.Error())
	}

	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	httpConfig := &serverconfigs.HTTPProtocolConfig{}
	if len(server.HttpJSON) > 0 {
		err = json.Unmarshal(server.HttpJSON, httpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		httpConfig.IsOn = true
	}

	httpConfig.Listen = addresses
	configData, err := json.Marshal(httpConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ServerRPC().UpdateServerHTTP(this.AdminContext(), &pb.UpdateServerHTTPRequest{
		ServerId: params.ServerId,
		Config:   configData,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
