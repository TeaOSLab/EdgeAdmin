package tcp

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	tcpConfig := &serverconfigs.TCPProtocolConfig{}
	if len(server.TcpJSON) > 0 {
		err := json.Unmarshal(server.TcpJSON, tcpConfig)
		if err != nil {
			this.ErrorPage(err)
		}
	} else {
		tcpConfig.IsOn = true
	}

	this.Data["serverType"] = server.Type
	this.Data["tcpConfig"] = tcpConfig

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId   int64
	ServerType string
	Addresses  string

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.ServerTCP_LogUpdateTCPSettings, params.ServerId)

	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}

	addresses := []*serverconfigs.NetworkAddressConfig{}
	err := json.Unmarshal([]byte(params.Addresses), &addresses)
	if err != nil {
		this.Fail("端口地址解析失败：" + err.Error())
	}

	tcpConfig := &serverconfigs.TCPProtocolConfig{}
	if len(server.TcpJSON) > 0 {
		err := json.Unmarshal(server.TcpJSON, tcpConfig)
		if err != nil {
			this.ErrorPage(err)
		}
	} else {
		tcpConfig.IsOn = true
	}
	tcpConfig.Listen = addresses

	configData, err := json.Marshal(tcpConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ServerRPC().UpdateServerTCP(this.AdminContext(), &pb.UpdateServerTCPRequest{
		ServerId: params.ServerId,
		TcpJSON:  configData,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
