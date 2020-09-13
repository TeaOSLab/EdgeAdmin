package tcp

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

// TCP设置
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("tcp")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	server, config, isOk := serverutils.FindServer(&this.ParentAction, params.ServerId)
	if !isOk {
		return
	}
	if config.TCP == nil {
		this.ErrorPage(errors.New("there is no tcp setting"))
		return
	}

	if config.TCP.Listen == nil {
		config.TCP.Listen = []*serverconfigs.NetworkAddressConfig{}
	}

	this.Data["serverType"] = server.Type
	this.Data["addresses"] = config.TCP.Listen

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId   int64
	ServerType string
	Addresses  string

	Must *actions.Must
}) {
	_, config, isOk := serverutils.FindServer(&this.ParentAction, params.ServerId)
	if !isOk {
		return
	}

	listen := []*serverconfigs.NetworkAddressConfig{}
	err := json.Unmarshal([]byte(params.Addresses), &listen)
	if err != nil {
		this.Fail("端口地址解析失败：" + err.Error())
	}

	if config.IsHTTP() {
		config.HTTP.Listen = listen
	} else if config.IsTCP() {
		config.TCP.Listen = listen
	} else if config.IsUnix() {
		config.Unix.Listen = listen
	} else if config.IsUDP() {
		config.UDP.Listen = listen
	}

	configData, err := config.AsJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ServerRPC().UpdateServerConfig(this.AdminContext(), &pb.UpdateServerConfigRequest{
		ServerId: params.ServerId,
		Config:   configData,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
