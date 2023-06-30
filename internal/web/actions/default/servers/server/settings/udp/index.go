package udp

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

// IndexAction UDP设置
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("udp")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	udpConfig := &serverconfigs.UDPProtocolConfig{}
	if len(server.UdpJSON) > 0 {
		err := json.Unmarshal(server.UdpJSON, udpConfig)
		if err != nil {
			this.ErrorPage(err)
		}
	} else {
		udpConfig.IsOn = true
	}

	this.Data["serverType"] = server.Type
	this.Data["udpConfig"] = udpConfig

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId   int64
	ServerType string
	Addresses  string

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.NSCluster_LogUpdateNSClusterSettingsUDP, params.ServerId)

	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}

	addresses := []*serverconfigs.NetworkAddressConfig{}
	err := json.Unmarshal([]byte(params.Addresses), &addresses)
	if err != nil {
		this.Fail("端口地址解析失败：" + err.Error())
	}

	udpConfig := &serverconfigs.UDPProtocolConfig{}
	if len(server.UdpJSON) > 0 {
		err := json.Unmarshal(server.UdpJSON, udpConfig)
		if err != nil {
			this.ErrorPage(err)
		}
	} else {
		udpConfig.IsOn = true
	}
	udpConfig.Listen = addresses

	configData, err := json.Marshal(udpConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ServerRPC().UpdateServerUDP(this.AdminContext(), &pb.UpdateServerUDPRequest{
		ServerId: params.ServerId,
		UdpJSON:  configData,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
