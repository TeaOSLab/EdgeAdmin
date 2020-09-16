package https

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
	this.SecondMenu("https")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	server, _, isOk := serverutils.FindServer(&this.ParentAction, params.ServerId)
	if !isOk {
		return
	}
	httpsConfig := &serverconfigs.HTTPSProtocolConfig{}
	if len(server.HttpsJSON) > 0 {
		err := json.Unmarshal(server.HttpsJSON, httpsConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		httpsConfig.IsOn = true
	}

	this.Data["serverType"] = server.Type
	this.Data["httpsConfig"] = maps.Map{
		"isOn":      httpsConfig.IsOn,
		"addresses": httpsConfig.Listen,
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

	server, _, isOk := serverutils.FindServer(&this.ParentAction, params.ServerId)
	if !isOk {
		return
	}
	httpsConfig := &serverconfigs.HTTPSProtocolConfig{}
	if len(server.HttpsJSON) > 0 {
		err = json.Unmarshal(server.HttpsJSON, httpsConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		httpsConfig.IsOn = true
	}

	httpsConfig.Listen = addresses
	configData, err := json.Marshal(httpsConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ServerRPC().UpdateServerHTTPS(this.AdminContext(), &pb.UpdateServerHTTPSRequest{
		ServerId: params.ServerId,
		Config:   configData,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
